import { PageShell, DataTable } from '../../components/PageShell'

const warehouses = [
  ['WH Jakarta Timur', 'Jakarta Timur, DKI Jakarta', '3,204 SKU', '78%'],
  ['WH Bandung', 'Bandung, Jawa Barat', '1,880 SKU', '54%'],
  ['WH Surabaya', 'Surabaya, Jawa Timur', '2,410 SKU', '61%'],
  ['WH Semarang', 'Semarang, Jawa Tengah', '960 SKU', '40%'],
]

export function WarehouseLocationsPage() {
  return (
    <PageShell
      title="Locations · Warehouse"
      description="Warehouse sites and how full each one is running."
      action="Add warehouse"
    >
      <DataTable
        columns={['Warehouse', 'Address', 'Items stored', 'Utilization']}
        rows={warehouses}
      />
    </PageShell>
  )
}
