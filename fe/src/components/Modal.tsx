import { useEffect } from 'react'

export function Modal(props: { open: boolean; title: string; onClose: () => void; children: React.ReactNode }) {
  useEffect(() => {
    if (!props.open) return
    const onKeyDown = (e: KeyboardEvent) => {
      if (e.key === 'Escape') props.onClose()
    }
    window.addEventListener('keydown', onKeyDown)
    return () => window.removeEventListener('keydown', onKeyDown)
  }, [props.open, props.onClose])

  if (!props.open) return null

  return (
    <div className="modalOverlay" role="dialog" aria-modal="true" aria-label={props.title}>
      <div className="modalOverlay__backdrop" onMouseDown={props.onClose} />
      <div className="modal">
        <div className="modal__header">
          <div className="modal__title">{props.title}</div>
          <button className="btn btn--ghost btn--icon" onClick={props.onClose} aria-label="Close">
            ✕
          </button>
        </div>
        <div className="modal__body">{props.children}</div>
      </div>
    </div>
  )
}


