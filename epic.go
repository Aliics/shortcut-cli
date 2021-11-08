package shortcut_api

import (
	"net/http"
	"strconv"
	"sync"
)

type Epic struct {
	Entity
	Completed bool      `json:"completed"`
	Started   bool      `json:"started"`
	Stats     EpicStats `json:"stats"`
}

type EpicStats struct {
	NumStoriesTotal   int `json:"num_stories_total"`
	NumStoriesStarted int `json:"num_stories_started"`
	NumStoriesDone    int `json:"num_stories_done"`
}

func (s Shortcut) ListEpics() ([]Epic, error) {
	req, err := http.NewRequest(
		"GET",
		s.url+"/epics",
		nil,
	)
	if err != nil {
		return nil, err
	}

	var result []Epic
	err = s.makeQuery(req, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}
func (s Shortcut) ListEpicStories(epic Epic) ([]Story, error) {
	req, err := http.NewRequest(
		"GET",
		s.url+"/epics/"+strconv.Itoa(epic.Id)+"/stories",
		nil,
	)
	if err != nil {
		return nil, err
	}

	var result []Story
	err = s.makeQuery(req, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s Shortcut) QueryEpics(predicate func(Epic) bool) ([]Epic, error) {
	allEpics, err := s.ListEpics()
	if err != nil {
		panic(err)
	}

	var epics []Epic
	for _, epic := range allEpics {
		if predicate(epic) {
			epics = append(epics, epic)
		}
	}

	return epics, nil
}

func (s Shortcut) ListStoriesForEpics(epics []Epic) ([]Story, error) {
	storiesCh := make(chan []Story, len(epics))

	go func() {
		storiesWg := sync.WaitGroup{}
		for i := 0; i < len(epics); i++ {
			storiesWg.Add(1)
			go func(i int) {
				story, err := s.ListEpicStories(epics[i])
				if err != nil {
					panic(err)
				}

				storiesCh <- story

				storiesWg.Done()
			}(i)
		}

		storiesWg.Wait()
		close(storiesCh)
	}()

	var stories []Story
	for ss := range storiesCh {
		stories = append(stories, ss...)
	}

	return stories, nil
}
