import { createResource } from 'solid-js'
import { ListAudioFileKeybindings, RemoveAudioFileKeybinding, SetAudioFileKeybinding } from '../wailsjs/go/main/App'

export const useAudioFileKeybindings = () => {
  // eslint-disable-next-line solid/reactivity
  const [data, { refetch }] = createResource(async () => {
    try {
      return await ListAudioFileKeybindings()
    }
    catch (err: unknown) {
      console.error(err)
    }
  }, { initialValue: {} })

  const set = async (audioFile: string, keybinding: string) => {
    await SetAudioFileKeybinding(audioFile, keybinding)
    await refetch()
  }

  const remove = async (audioFile: string) => {
    await RemoveAudioFileKeybinding(audioFile)
    await refetch()
  }

  return {
    audioFileKeybindings: data,
    refetchAudioFileKeybindings: refetch,
    setAudioFileKeybinding: set,
    removeAudioFileKeybinding: remove,
  }
}
