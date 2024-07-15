import { createResource } from 'solid-js'
import { AddAudioFile, ListAudioFiles, PlayAudioFile, RemoveAudioFile, StopAudioFile } from '../wailsjs/go/main/App'

export const useAudioFiles = () => {
  // eslint-disable-next-line solid/reactivity
  const [data, { refetch }] = createResource(async () => {
    try {
      return await ListAudioFiles()
    }
    catch (err) {
      console.error(err)
    }
  }, { initialValue: [] })

  const add = async (file: string) => {
    await AddAudioFile(file)
    await refetch()
  }

  const remove = async (file: string) => {
    await RemoveAudioFile(file)
    await refetch()
  }

  const play = async (file: string) => {
    await PlayAudioFile(file)
  }

  const stop = async (file: string) => {
    await StopAudioFile(file)
  }

  return {
    audioFiles: data,
    addAudioFile: add,
    removeAudioFile: remove,
    playAudioFile: play,
    stopAudioFile: stop,
  }
}
