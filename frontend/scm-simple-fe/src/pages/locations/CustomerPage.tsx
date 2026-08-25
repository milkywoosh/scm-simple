import { PageShell, DataTable } from '../../components/PageShell'

const customerSites = [
  ['PT Maju Bersama', 'Jakarta Timur', '1x ONT Router Type C', '12 Aug 2026'],
  ['CV Cahaya Abadi', 'Bandung', '1x ONT Router Type C', '09 Aug 2026'],
  ['Toko Sumber Rejeki', 'Surabaya', '10m Fiber Patch Cord', '02 Aug 2026'],
]

export function CustomerLocationsPage() {
  return (
    <PageShell
      title="Locations · Customer"
      description="Equipment installed and currently in service at a customer site."
    >
      <DataTable
        columns={['Customer', 'Region', 'Installed equipment', 'Installed on']}
        rows={customerSites}
      />
    </PageShell>
  )
}
