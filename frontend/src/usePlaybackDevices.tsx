import { createResource } from 'solid-js'
import { ListPlaybackDevices } from '../wailsjs/go/main/App'

export const usePlaybackDevices = () => {
  // eslint-disable-next-line solid/reactivity
  const [data, { refetch }] = createResource(async () => {
    try {
      return await ListPlaybackDevices()
    }
    catch (e) {
      console.error(e)
    }
  }, { initialValue: [] })

  return {
    playbackDevices: data,
    refetchPlaybackDevices: refetch,
  }
}
