import { PageShell, DataTable } from '../../components/PageShell'

const technicianStock = [
  ['Rian Pratama', 'Jakarta Timur', '2x ONT Router, 30m Drop Cable'],
  ['Sari Wulandari', 'Bandung', '4x Wall Mount Bracket'],
  ['Dedi Kurniawan', 'Jakarta Timur', '5x Wall Mount Bracket'],
]

export function TechnicianLocationsPage() {
  return (
    <PageShell
      title="Locations · Technician"
      description="Field stock currently held in a technician's kit, not yet installed."
    >
      <DataTable
        columns={['Technician', 'Region', 'Stock on hand']}
        rows={technicianStock}
      />
    </PageShell>
  )
}
