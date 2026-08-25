import { PageShell, DataTable } from '../../components/PageShell'

const suppliers = [
  ['Nusantara Fiber Optics', 'purchasing@nusantarafo.id', '5 days', '3'],
  ['PT Sinar Kabel Jaya', 'sales@sinarkabel.co.id', '7 days', '1'],
  ['Global Telecom Supply', 'orders@gts-supply.com', '12 days', '2'],
]

export function SuppliersPage() {
  return (
    <PageShell
      title="Suppliers"
      description="Vendors this project sources materials and equipment from."
      action="Add supplier"
    >
      <DataTable
        columns={['Supplier', 'Contact', 'Lead time', 'Open POs']}
        rows={suppliers}
      />
    </PageShell>
  )
}
