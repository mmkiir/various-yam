import { type Component, For } from 'solid-js'
import { useCaptureDeviceID } from './useCaptureDeviceID'
import { useCaptureDevices } from './useCaptureDevices'
import { usePlaybackDeviceID } from './usePlaybackDeviceID'
import { usePlaybackDevices } from './usePlaybackDevices'

const App: Component = () => {
  const { captureDeviceID, setCaptureDeviceID } = useCaptureDeviceID()
  const { captureDevices, refetchCaptureDevices } = useCaptureDevices()
  const { playbackDeviceID, setPlaybackDeviceID } = usePlaybackDeviceID()
  const { playbackDevices, refetchPlaybackDevices } = usePlaybackDevices()

  const handleCaptureDeviceIDChange = (event: Event & {
    currentTarget: HTMLSelectElement
    target: HTMLSelectElement
  }) => {
    setCaptureDeviceID((event.currentTarget).value).catch((err: unknown) => {
      console.error(err)
    })
  }

  const handlePlaybackDeviceIDChange = (event: Event & {
    currentTarget: HTMLSelectElement
    target: HTMLSelectElement
  }) => {
    setPlaybackDeviceID((event.currentTarget).value).catch((err: unknown) => {
      console.error(err)
    })
  }

  const handleCaptureDevicesFocus = () => {
    Promise.resolve(refetchCaptureDevices()).catch((err: unknown) => {
      console.error(err)
    })
  }

  const handlePlaybackDevicesFocus = () => {
    Promise.resolve(refetchPlaybackDevices()).catch((err: unknown) => {
      console.error(err)
    })
  }

  return (
    <>
      <main class="container-fluid">
        <select
          onChange={handleCaptureDeviceIDChange}
          onFocus={handleCaptureDevicesFocus}
        >
          <For each={captureDevices()}>
            {device => (
              <option
                selected={device.deviceId === captureDeviceID()}
                value={device.deviceId}
              >
                {device.label}
              </option>
            )}
          </For>
        </select>
        <select
          onChange={handlePlaybackDeviceIDChange}
          onFocus={handlePlaybackDevicesFocus}
        >
          <For each={playbackDevices()}>
            {device => (
              <option
                selected={device.deviceId === playbackDeviceID()}
                value={device.deviceId}
              >
                {device.label}
              </option>
            )}
          </For>
        </select>
      </main>
    </>
  )
}

export default App
