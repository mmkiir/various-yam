import { createResource } from 'solid-js'
import { GetCaptureDeviceID, SetCaptureDeviceID } from '../wailsjs/go/main/App'

export const useCaptureDeviceID = () => {
  // eslint-disable-next-line solid/reactivity
  const [data, { refetch }] = createResource(async () => {
    try {
      return await GetCaptureDeviceID()
    }
    catch (e) {
      console.error(e)
    }
  }, { initialValue: '' })

  return {
    captureDeviceID: data,
    refetchCaptureDeviceID: refetch,
    setCaptureDeviceID: async (id: string) => {
      await SetCaptureDeviceID(id)
      await refetch()
    },
  }
}
