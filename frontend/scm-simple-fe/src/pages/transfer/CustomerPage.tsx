import { PageShell, DataTable } from '../../components/PageShell'

const deliveries = [
  ['TRF-4108', 'Rian Pratama', 'PT Maju Bersama', '1x ONT Router Type C', 'Installed'],
  ['TRF-4102', 'Sari Wulandari', 'CV Cahaya Abadi', '1x ONT Router Type C', 'Installed'],
  ['TRF-4097', 'Dedi Kurniawan', 'Toko Sumber Rejeki', '10m Fiber Patch Cord', 'Installed'],
]

export function CustomerTransferPage() {
  return (
    <PageShell
      title="Transfer · Customer"
      description="Equipment handed off from a technician and installed at a customer site."
    >
      <DataTable
        columns={['Ref', 'Technician', 'Customer', 'Contents', 'Status']}
        rows={deliveries}
      />
    </PageShell>
  )
}
