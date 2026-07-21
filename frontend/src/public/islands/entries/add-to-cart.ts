import { z } from 'zod'

import mountIsland from '../../../mount'
import AddToCart from '../products/AddToCart'

const el = document.getElementById('add-to-cart')
if (!el) {
  throw new Error('add-to-cart not found')
}

const res = z.uuid().safeParse(el.dataset.productId)
if (!res.success) {
  throw new Error('invalid or missing product id')
}

mountIsland({
  el,
  Component: AddToCart,
  props: { productId: res.data },
  hydrate: true,
})
