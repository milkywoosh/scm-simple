import { useState } from 'react'
import { PageShell, DataTable } from '../../components/PageShell'

const allItems = [
  ['SKU-10231', 'Fiber Patch Cord 10m', 'Cable', '842'],
  ['SKU-10442', 'ONT Router Type C', 'Equipment', '156'],
  ['SKU-10118', 'Splice Closure 24F', 'Enclosure', '389'],
  ['SKU-10509', 'Drop Cable 2-Core', 'Cable', '1204'],
  ['SKU-10087', 'Fusion Splicer Blade', 'Tooling', '22'],
]

export function SearchPage() {
  const [query, setQuery] = useState('')

  const results = allItems.filter((row) =>
    row.some((cell) => cell.toLowerCase().includes(query.toLowerCase())),
  )

  return (
    <PageShell
      title="Search"
      description="Find items across every warehouse, technician kit, and customer site by SKU, name, or category."
    >
      <input
        value={query}
        onChange={(e) => setQuery(e.target.value)}
        placeholder="Search by SKU or item name..."
        className="mb-4 w-full rounded-md border border-line bg-auto px-4 py-2.5 text-sm text-ink placeholder:text-ink-soft focus:border-rail focus:outline-none"
      />
      <DataTable columns={['SKU', 'Item', 'Category', 'On hand']} rows={results} />
      {results.length === 0 && (
        <p className="mt-4 text-center text-sm text-ink-soft">
          No items match "{query}". Try a different SKU or keyword.
        </p>
      )}
    </PageShell>
  )
}
