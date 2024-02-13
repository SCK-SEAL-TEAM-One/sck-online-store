'use client'

import Button from '@/components/button/button'
import InputField from '@/components/input-field'
import Tag from '@/components/tag'

// ----------------------------------------------------------------------

const DiscountForm = () => {
  return (
    <div className="border-b pb-5">
      <div className="-mx-2 flex items-end justify-between">
        <div className="flex-grow px-2 lg:max-w-sm">
          <InputField
            id="discount-form-discount-code"
            label="discount"
            type="text"
            placeholder="XXXXXX"
          />
        </div>
        <div className="px-2">
          <Button id="discount-apply-btn" size="sm">
            APPLY
          </Button>
        </div>
      </div>

      <div className="mt-2">
        <Tag name="SAVE20" onRemove={() => {}} />
        <Tag name="FREESHIP" onRemove={() => {}} />
      </div>
    </div>
  )
}

export default DiscountForm
