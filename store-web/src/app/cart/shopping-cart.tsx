'use client'

import ButtonContiueShopping from '@/app/cart/components/button-continue'
import ProductList from '@/app/cart/components/product-list'
import ShoppingHeader from '@/app/cart/components/shopping-header'
import SubTotal from '@/app/cart/components/sub-total'
import ButtonLink from '@/components/button/button-link'
import Text from '@/components/typography/text'
import useOrderStore from '@/hooks/use-order-store'
import { Dialog, Transition } from '@headlessui/react'
import { Fragment, useEffect } from 'react'

// ----------------------------------------------------------------------

type ShoppingCartViewProps = {
  openShoppingCart: boolean
  setOpenShoppingCart: (open: boolean) => void
}

const ShoppingCartView = ({
  openShoppingCart,
  setOpenShoppingCart
}: ShoppingCartViewProps) => {
  const { cart, summary, getProductListInCart } = useOrderStore()

  useEffect(() => {
    getProductListInCart()
  }, [getProductListInCart])

  return (
    <Transition.Root show={openShoppingCart} as={Fragment}>
      <Dialog as="div" className="relative z-10" onClose={setOpenShoppingCart}>
        <Transition.Child
          as={Fragment}
          enter="ease-in-out duration-500"
          enterFrom="opacity-0"
          enterTo="opacity-100"
          leave="ease-in-out duration-500"
          leaveFrom="opacity-100"
          leaveTo="opacity-0"
        >
          <div
            id="shopping-cart-overlay"
            className="fixed inset-0 bg-gray-500 bg-opacity-75 transition-opacity"
          />
        </Transition.Child>

        <div className="fixed inset-0 overflow-hidden">
          <div className="absolute inset-0 overflow-hidden">
            <div className="pointer-events-none fixed inset-y-0 right-0 flex max-w-full pl-10">
              <Transition.Child
                as={Fragment}
                enter="transform transition ease-in-out duration-500 sm:duration-700"
                enterFrom="translate-x-full"
                enterTo="translate-x-0"
                leave="transform transition ease-in-out duration-500 sm:duration-700"
                leaveFrom="translate-x-0"
                leaveTo="translate-x-full"
              >
                <Dialog.Panel className="pointer-events-auto w-screen max-w-xl">
                  <div className="flex h-full flex-col overflow-y-scroll bg-white shadow-xl">
                    <div className="flex-1 overflow-y-auto px-4 py-6 sm:px-6">
                      <ShoppingHeader
                        id="shopping-cart-header"
                        closeShoppingCart={() => setOpenShoppingCart(false)}
                      />

                      <ProductList list={cart} />
                    </div>

                    <div className="border-t border-gray-200 px-4 py-6 sm:px-6">
                      <SubTotal total={summary.total_price_thb} />

                      <Text size="sm" className="mt-0.5 text-gray-500">
                        Shipping and taxes calculated at checkout.
                      </Text>

                      <ButtonLink
                        id="shopping-cart-checkout-btn"
                        href="/checkout"
                        className="mt-6"
                      >
                        Checkout
                      </ButtonLink>

                      <div className="mt-6 flex justify-center flex-col items-center text-center text-sm text-gray-500">
                        <Text className="mb-3">or</Text>
                        <ButtonContiueShopping
                          id="shopping-cart-continue-link"
                          onClick={() => setOpenShoppingCart(false)}
                        />
                      </div>
                    </div>
                  </div>
                </Dialog.Panel>
              </Transition.Child>
            </div>
          </div>
        </div>
      </Dialog>
    </Transition.Root>
  )
}

export default ShoppingCartView
