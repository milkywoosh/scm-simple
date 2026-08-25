import { PageShell, DataTable } from '../../components/PageShell'

const products = [
  ['SKU-10231', 'Fiber Patch Cord 10m', 'Cable', 'unit', '842'],
  ['SKU-10442', 'ONT Router Type C', 'Equipment', 'unit', '156'],
  ['SKU-10118', 'Splice Closure 24F', 'Enclosure', 'unit', '389'],
  ['SKU-10509', 'Drop Cable 2-Core', 'Cable', 'meter', '1204'],
  ['SKU-10087', 'Fusion Splicer Blade', 'Tooling', 'unit', '22'],
  ['SKU-10375', 'Wall Mount Bracket', 'Hardware', 'unit', '640'],
]

export function ProductsPage() {
  return (
    <PageShell
      title="Products"
      description="The master catalog of items tracked across the supply chain."
      action="Add product"
    >
      <DataTable
        columns={['SKU', 'Name', 'Category', 'Unit', 'On hand']}
        rows={products}
      />
    </PageShell>
  )
}
