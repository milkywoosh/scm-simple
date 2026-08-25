import { PageShell, DataTable } from '../../components/PageShell'

const transfers = [
  ['TRF-2208', 'WH Jakarta Timur', 'WH Bandung', '120x Drop Cable', 'In transit'],
  ['TRF-2207', 'WH Bandung', 'WH Surabaya', '40x ONT Router', 'Delivered'],
  ['TRF-2201', 'WH Jakarta Timur', 'WH Semarang', '18x Splicer Blade', 'Delivered'],
]

export function WarehouseTransferPage() {
  return (
    <PageShell
      title="Transfer · Warehouse"
      description="Stock moves between warehouses."
      action="New transfer"
    >
      <DataTable
        columns={['Ref', 'From', 'To', 'Contents', 'Status']}
        rows={transfers}
      />
    </PageShell>
  )
}
