import { createResource } from 'solid-js'
import { GetPlaybackDeviceID, SetPlaybackDeviceID } from '../wailsjs/go/main/App'

export const usePlaybackDeviceID = () => {
  // eslint-disable-next-line solid/reactivity
  const [data, { refetch }] = createResource(async () => {
    try {
      return await GetPlaybackDeviceID()
    }
    catch (e) {
      console.error(e)
    }
  }, { initialValue: '' })

  const set = async (id: string) => {
    await SetPlaybackDeviceID(id)
    await refetch()
  }

  return {
    playbackDeviceID: data,
    refetchPlaybackDeviceID: refetch,
    setPlaybackDeviceID: set,
  }
}
