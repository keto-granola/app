import { useState } from 'react'

interface AddToCartProps {
  productId: string
}

function AddToCart({ productId }: AddToCartProps) {
  const [quantity, setQuantity] = useState<number>(1)

  const addToCart = (productId: string, quantity: number) => {
    // TODO: add to cart
    console.warn(productId, quantity)
  }

  return (
    <div>
      <button onClick={() => setQuantity(q => q - 1)}>-</button>
      <span>{quantity}</span>
      <button onClick={() => setQuantity(q => q + 1)}>+</button>
      <button onClick={() => addToCart(productId, quantity)}>Add to cart</button>
    </div>
  )
}

export default AddToCart
