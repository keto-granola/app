import type { ComponentType } from 'react'
import { createRoot, hydrateRoot } from 'react-dom/client'

import './index.css'

interface MountOpts<P> {
  el: HTMLElement
  Component: ComponentType<P>
  props: P
  hydrate?: boolean
}

function mountIsland<P extends object>({ el, Component, props, hydrate = false }: MountOpts<P>) {
  if (hydrate) {
    hydrateRoot(el, <Component {...props} />)
  } else {
    createRoot(el).render(<Component {...props} />)
  }
}

export default mountIsland
