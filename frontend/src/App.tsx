import { type Component, For } from 'solid-js'
import { OpenMultipleFilesDialog } from '../wailsjs/go/main/App'
import { main } from '../wailsjs/go/models'
import { useAudioFiles } from './useAudioFiles'
import { useCaptureDeviceID } from './useCaptureDeviceID'
import { useCaptureDevices } from './useCaptureDevices'
import { usePlaybackDeviceID } from './usePlaybackDeviceID'
import { usePlaybackDevices } from './usePlaybackDevices'

const App: Component = () => {
  const { audioFiles, addAudioFile, removeAudioFile, playAudioFile, stopAudioFile } = useAudioFiles()
  const { captureDeviceID, setCaptureDeviceID } = useCaptureDeviceID()
  const { captureDevices, refetchCaptureDevices } = useCaptureDevices()
  const { playbackDeviceID, setPlaybackDeviceID } = usePlaybackDeviceID()
  const { playbackDevices, refetchPlaybackDevices } = usePlaybackDevices()

  const handleCaptureDeviceIDChange = (event: Event & { currentTarget: HTMLSelectElement, target: HTMLSelectElement }) => {
    setCaptureDeviceID(event.currentTarget.value).catch((err: unknown) => {
      console.error(err)
    })
  }

  const handlePlaybackDeviceIDChange = (event: Event & { currentTarget: HTMLSelectElement, target: HTMLSelectElement }) => {
    setPlaybackDeviceID(event.currentTarget.value).catch((err: unknown) => {
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

  const handleAction = (action: () => Promise<void>) => () => {
    action().catch((err: unknown) => {
      console.error(err)
    })
  }

  const handleOpenMultipleFilesDialog = () => {
    OpenMultipleFilesDialog({
      title: 'Select audio files',
      filters: [{ displayName: 'Audio files', pattern: '*.mp3;*.wav;*.ogg' }],
    } as main.OpenDialogOptions)
      .then((files) => {
        files.forEach((file) => {
          addAudioFile(file).catch((err: unknown) => {
            console.error(err)
          })
        })
      })
      .catch((err: unknown) => {
        console.error(err)
      })
  }

  return (
    <>
      <header
        class="container-fluid"
        style={{
          '--pico-form-element-spacing-horizontal': '0.5rem',
          '--pico-form-element-spacing-vertical': '0.375rem',
        }}
      >
        <nav>
          <ul>
            <li>
              <button
                class="outline"
                onClick={handleOpenMultipleFilesDialog}
              >
                📂
              </button>
            </li>
          </ul>
        </nav>
      </header>
      <main
        class="container-fluid"
        style={{
          '--pico-form-element-spacing-horizontal': '0.5rem',
          '--pico-form-element-spacing-vertical': '0.375rem',
        }}
      >
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
        <ul>
          <For each={audioFiles()}>
            {audioFile => (
              <li style={{
                display: 'flex',
                gap: 'var(--pico-spacing)',
              }}
              >
                <span style={{ flex: 1 }}>{audioFile}</span>
                <For each={[
                  ['▶️', () => playAudioFile(audioFile)],
                  ['⏹️', () => stopAudioFile(audioFile)],
                  ['🗑️', () => removeAudioFile(audioFile)],
                ] as const}
                >
                  {([emoji, action]) => (
                    <button class="outline" onClick={handleAction(action)}>
                      {emoji}
                    </button>
                  )}
                </For>
              </li>
            )}
          </For>
        </ul>
      </main>
    </>
  )
}

export default App
