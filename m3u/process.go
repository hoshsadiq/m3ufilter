package m3u

import (
	"fmt"
	"github.com/hoshsadiq/m3ufilter/config"
	"github.com/hoshsadiq/m3ufilter/logger"
	"github.com/hoshsadiq/m3ufilter/m3u/csv"
	"github.com/hoshsadiq/m3ufilter/m3u/helper"
	"github.com/hoshsadiq/m3ufilter/m3u/xmltv"
	"github.com/hoshsadiq/m3ufilter/net"
	"sort"
	"strings"
)

var log = logger.Get()

func processProvider(conf *config.Config, provider *config.Provider, epg *xmltv.XMLTV, generatingCsv bool) Streams {
	resp, err := net.GetUri(provider.Uri)
	if err != nil {
		log.Errorf("Could not retrieve file from %s, err = %v", provider.Uri, err)
		return nil
	}
	defer helper.Close(resp.Body, fmt.Sprintf("for provider %s", provider.Uri))()

	var csvData map[string]*csv.StreamData
	if !generatingCsv {
		csvData, err = csv.GetCsvMapping(provider.Csv)
		if err != nil {
			log.Errorf("Could not retrieve file from %s, err = %v", provider.Csv, err)
			return nil
		}
	}

	pl, err := decode(conf, resp.Body, csvData, provider.CheckStreams, epg)
	if err != nil {
		log.Errorf("could not decode playlist from provider %s, err = %v", provider.Uri, err)
		return nil
	}
	return pl
}

func ProcessConfig(conf *config.Config, generatingCsv bool) (streams Streams, epg *xmltv.XMLTV, allFailed bool) {
	errs := 0
	var err error

	if !generatingCsv {
		epg, err = getEpg(conf.EpgProviders)
		if err != nil {
			log.Errorf("Could not parse EPG, skipping all EPG related tasks; err=%v", err)
		}
	} else {
		epg = nil
	}

	streams = Streams{}
	// todo we can do each provider in its own coroutine, then converged at the end.
	//   furthermore, each line can be done in its own coroutine as well.
	var pl Streams
	for _, provider := range conf.Providers {
		pl = processProvider(conf, provider, epg, generatingCsv)

		if pl == nil {
			errs++
		} else {
			streams = append(streams, pl...)
		}
	}

	if generatingCsv {
		return streams, nil, false
	}

	sort.Sort(streams)

	if epg != nil {
		var streamIds = make(map[string]bool, len(streams))
		for _, stream := range streams {
			if stream.Id != "" {
				streamIds[strings.ToLower(stream.Id)] = true
			}
		}

		var newProgrammes = make([]*xmltv.Programme, 0, len(epg.Programmes))
		for _, programme := range epg.Programmes {
			programme.Channel = strings.ToLower(programme.Channel)
			if _, ok := streamIds[programme.Channel]; ok {
				newProgrammes = append(newProgrammes, programme)
			}
		}
		epg.Programmes = newProgrammes

		var addChannel bool
		var newEpgChannels = make([]*xmltv.Channel, 0, len(epg.Channels))
		for _, epgChannel := range epg.Channels {
			epgChannel.ID = strings.ToLower(epgChannel.ID)

			if _, ok := streamIds[epgChannel.ID]; ok {
				addChannel = true
				for _, newEpgChannel := range newEpgChannels {
					if newEpgChannel.ID == epgChannel.ID {
						addChannel = false
						newEpgChannel.DisplayNames = append(newEpgChannel.DisplayNames, epgChannel.DisplayNames...)
						break
					}
				}
				if addChannel {
					newEpgChannels = append(newEpgChannels, epgChannel)
				}
			}
		}
		epg.Channels = newEpgChannels
		setEpgInfo(epg)
	}

	return streams, epg, len(conf.Providers) == errs
}

func setEpgInfo(epg *xmltv.XMLTV) {
	epg.SetGenerator(config.EpgGeneratorName(), config.EpgGeneratorUrl())
}

func getEpg(providers []*config.EpgProvider) (*xmltv.XMLTV, error) {
	var epgs = make([]xmltv.XMLTV, len(providers))
	var newEpg xmltv.XMLTV
	totalChannels := 0
	totalProgrammes := 0

	for i, provider := range providers {
		resp, err := net.GetUri(provider.Uri)
		if err != nil {
			return nil, err
		}
		defer func() {
			err := resp.Body.Close()
			if err != nil {
				log.Errorf("could not close request body for provider %s, err = %v", provider.Uri, err)
			}
		}()

		newEpg = xmltv.XMLTV{}
		err = xmltv.Load(resp.Body, &newEpg)
		if err != nil {
			return nil, err
		}
		applyEpgIdRenames(&newEpg, provider.ChannelIdRenames)
		epgs[i] = newEpg
		totalChannels += len(newEpg.Channels)
		totalProgrammes += len(newEpg.Programmes)
	}

	allChannels := make([]*xmltv.Channel, 0, totalChannels)
	allProgrammes := make([]*xmltv.Programme, 0, totalChannels)
	for _, epg := range epgs {
		allChannels = append(allChannels, epg.Channels...)
		allProgrammes = append(allProgrammes, epg.Programmes...)
	}

	var channels = make(map[string]*xmltv.Channel, len(allChannels))
	var nameIdMapping = make(map[string]string)

	for _, c := range allChannels {
		channel, ok := channels[c.ID]
		var found = false
		if !ok {
			for _, dpname := range c.DisplayNames {
				channelId, ok := nameIdMapping[dpname.Value]
				if ok {
					channel, ok = channels[channelId]
					nameIdMapping[c.ID] = channelId
					found = true
					break
				}
			}

			if !found {
				channels[c.ID] = c
				for _, dpname := range c.DisplayNames {
					nameIdMapping[dpname.Value] = c.ID
				}
				continue
			}
		}
		for _, left := range c.DisplayNames {
			found = false
			for _, right := range channel.DisplayNames {
				if right.Value == left.Value {
					found = true
					break
				}
			}
			if !found {
				channel.DisplayNames = append(channel.DisplayNames, left)
			}

			nameIdMapping[left.Value] = c.ID
		}
	}

	log.Info("Finished loading EPG")

	return &xmltv.XMLTV{Programmes: allProgrammes, Channels: allChannels}, nil
}

func applyEpgIdRenames(epg *xmltv.XMLTV, renames []config.ChannelIdRename) {
	for _, rename := range renames {
		for _, chann := range epg.Channels {
			if chann.ID == rename.From {
				chann.ID = rename.To
			}
		}

		for _, programme := range epg.Programmes {
			if programme.Channel == rename.From {
				programme.Channel = rename.To
			}
		}
	}
}
