import { createResource } from 'solid-js'
import { ListCaptureDevices } from '../wailsjs/go/main/App'

export const useCaptureDevices = () => {
  // eslint-disable-next-line solid/reactivity
  const [data, { refetch }] = createResource(async () => {
    try {
      return await ListCaptureDevices()
    }
    catch (e) {
      console.error(e)
    }
  }, { initialValue: [] })

  return {
    captureDevices: data,
    refetchCaptureDevices: refetch,
  }
}
