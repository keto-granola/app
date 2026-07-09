import { mountIsland } from '../../mount'
import { AddToCart } from '../components/products/AddToCart'
import '../../index.css'

const el = document.getElementById('add-to-cart')
if (el) {
  const productId = el.dataset.productId
  if (!productId) {
    throw new Error('missing product id')
  }

  mountIsland({
    el,
    Component: AddToCart,
    props: { productId },
  })
}
