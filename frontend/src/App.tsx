import { type Component, createSignal, For, Show } from 'solid-js'
import { OpenMultipleFilesDialog } from '../wailsjs/go/main/App'
import { main } from '../wailsjs/go/models'
import { useAudioFileKeybindings } from './useAudioFileKeybindings'
import { useAudioFiles } from './useAudioFiles'
import { useCaptureDeviceID } from './useCaptureDeviceID'
import { useCaptureDevices } from './useCaptureDevices'
import { usePlaybackDeviceID } from './usePlaybackDeviceID'
import { usePlaybackDevices } from './usePlaybackDevices'

const App: Component = () => {
  const { audioFileKeybindings, setAudioFileKeybinding, removeAudioFileKeybinding } = useAudioFileKeybindings()
  const { audioFiles, addAudioFile, removeAudioFile, playAudioFile, stopAudioFile } = useAudioFiles()
  const { captureDeviceID, setCaptureDeviceID } = useCaptureDeviceID()
  const { captureDevices, refetchCaptureDevices } = useCaptureDevices()
  const { playbackDeviceID, setPlaybackDeviceID } = usePlaybackDeviceID()
  const { playbackDevices, refetchPlaybackDevices } = usePlaybackDevices()

  const handleCaptureDeviceIDChange = async (event: Event & { currentTarget: HTMLSelectElement, target: HTMLSelectElement }) => {
    await setCaptureDeviceID(event.currentTarget.value)
  }

  const handlePlaybackDeviceIDChange = async (event: Event & { currentTarget: HTMLSelectElement, target: HTMLSelectElement }) => {
    await setPlaybackDeviceID(event.currentTarget.value)
  }

  const handleCaptureDevicesFocus = async () => {
    await refetchCaptureDevices()
  }

  const handlePlaybackDevicesFocus = async () => {
    await refetchPlaybackDevices()
  }

  const handleOpenMultipleFilesDialog = async () => {
    const files = await OpenMultipleFilesDialog({
      title: 'Select audio files',
      filters: [{ displayName: 'Audio files', pattern: '*.mp3;*.wav;*.ogg' }],
    } as main.OpenDialogOptions)

    files.forEach((file) => {
      addAudioFile(file).catch((err: unknown) => {
        console.error(err)
      })
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
                onClick={() => {
                  handleOpenMultipleFilesDialog().catch((err: unknown) => {
                    console.error(err)
                  })
                }}
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
          onChange={(event) => {
            handleCaptureDeviceIDChange(event).catch((err: unknown) => {
              console.error(err)
            })
          }}
          onFocus={() => {
            handleCaptureDevicesFocus().catch((err: unknown) => {
              console.error(err)
            })
          }}
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
          onChange={(event) => {
            handlePlaybackDeviceIDChange(event).catch((err: unknown) => {
              console.error(err)
            })
          }}
          onFocus={() => {
            handlePlaybackDevicesFocus().catch((err: unknown) => {
              console.error(err)
            })
          }}
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
            {(audioFile) => {
              const [heldKeys, setHeldKeys] = createSignal<string[]>([])
              const [isRecording, setIsRecording] = createSignal(false)

              let dialog: HTMLDialogElement | undefined

              const handleKeyDown = (event: KeyboardEvent) => {
                event.preventDefault()
                if (!isRecording()) {
                  setHeldKeys([])
                  setIsRecording(true)
                }
                if (!heldKeys().includes(event.key)) {
                  setHeldKeys([...heldKeys(), event.key])
                }
              }

              const handleKeyUp = () => {
                if (isRecording()) {
                  setIsRecording(false)
                }
              }

              const handleSave = async () => {
                if (heldKeys().length === 0) {
                  return
                }

                await setAudioFileKeybinding(audioFile, heldKeys().join(' + '))
                setHeldKeys([])
                setIsRecording(false)
              }

              return (
                <li style={{
                  display: 'flex',
                  gap: 'calc(var(--pico-spacing) / 2)',
                }}
                >
                  <span style={{ flex: 1 }}>
                    {audioFile}
                  </span>
                  <Show when={audioFileKeybindings()?.[audioFile] !== undefined}>
                    <kbd style={{ 'align-self': 'center' }}>
                      {audioFileKeybindings()?.[audioFile] ?? ''}
                    </kbd>
                  </Show>
                  <button
                    class="outline"
                    onClick={() => {
                      dialog?.show()
                    }}
                  >
                    ⌨️
                  </button>
                  <button
                    class="outline"
                    onClick={() => {
                      playAudioFile(audioFile).catch((err: unknown) => {
                        console.error(err)
                      })
                    }}
                  >
                    ▶️
                  </button>
                  <button
                    class="outline"
                    onClick={() => {
                      stopAudioFile(audioFile).catch((err: unknown) => {
                        console.error(err)
                      })
                    }}
                  >
                    ⏹️
                  </button>
                  <button
                    class="outline"
                    onClick={() => {
                      removeAudioFile(audioFile).catch((err: unknown) => {
                        console.error(err)
                      })
                    }}
                  >
                    🗑️
                  </button>
                  <dialog ref={dialog}>
                    <article>
                      <header>
                        {audioFile}
                      </header>
                      <fieldset role="group">
                        <input
                          type="text"
                          value={
                            heldKeys().length > 0
                              ? heldKeys().join(' + ')
                              : audioFileKeybindings()?.[audioFile] ?? ''
                          }
                          onKeyDown={handleKeyDown}
                          onKeyUp={handleKeyUp}
                        />
                        <button
                          onClick={() => {
                            setHeldKeys([])
                            setIsRecording(false)
                            removeAudioFileKeybinding(audioFile).catch((err: unknown) => {
                              console.error(err)
                            })
                          }}
                        >
                          ❌
                        </button>
                      </fieldset>
                      <footer>
                        <button
                          onClick={() => {
                            dialog?.close()
                          }}
                        >
                          Close
                        </button>
                        <button
                          onClick={() => {
                            handleSave().catch((err: unknown) => {
                              console.error(err)
                            })
                            dialog?.close()
                          }}
                        >
                          Save
                        </button>
                      </footer>
                    </article>
                  </dialog>
                </li>
              )
            }}
          </For>
        </ul>
      </main>
    </>
  )
}

export default App
