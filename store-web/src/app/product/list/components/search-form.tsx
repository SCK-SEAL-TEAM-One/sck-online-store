'use client'

import Button from '@/components/button/button'
import { MagnifyingGlassIcon } from '@heroicons/react/24/outline'

// ----------------------------------------------------------------------

type SearchFormProps = {
  keyword: string
  onChangeSearchKeyword: (e: React.ChangeEvent<HTMLInputElement>) => void
  onSubmitSearch: (e: React.FormEvent<HTMLFormElement>) => void
}

const SearchForm = ({
  keyword,
  onChangeSearchKeyword,
  onSubmitSearch
}: SearchFormProps) => {
  return (
    <form className="flex items-center my-5" onSubmit={onSubmitSearch}>
      <div className="relative w-full">
        <div id="search-product-icon" className="flex absolute inset-y-0 left-0 items-center pl-3 pointer-events-none">
          <MagnifyingGlassIcon className="w-5 h-5 text-gray-500 dark:text-gray-400" />
        </div>
        <input
          type="text"
          id="search-product-input"
          className="bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-indigo-600 focus:border-indigo-600 block w-full pl-10 py-3"
          placeholder="Search product name..."
          onChange={onChangeSearchKeyword}
          value={keyword}
        />
      </div>
      <Button id="search-product-btn" className="inline-flex items-center py-2 px-3 ml-2 text-sm font-medium">
        <MagnifyingGlassIcon className="mr-2 -ml-1 w-5 h-5" />
        Search
      </Button>
    </form>
  )
}

export default SearchForm
