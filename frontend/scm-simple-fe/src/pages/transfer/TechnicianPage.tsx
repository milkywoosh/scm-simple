import { PageShell, DataTable } from '../../components/PageShell'

const checkouts = [
  ['TRF-3312', 'WH Jakarta Timur', 'Rian Pratama', '2x ONT Router, 30m Drop Cable', 'Checked out'],
  ['TRF-3309', 'WH Bandung', 'Sari Wulandari', '1x Splice Closure 24F', 'Returned'],
  ['TRF-3301', 'WH Jakarta Timur', 'Dedi Kurniawan', '5x Wall Mount Bracket', 'Checked out'],
]

export function TechnicianTransferPage() {
  return (
    <PageShell
      title="Transfer · Technician"
      description="Items checked out from a warehouse to a field technician's kit."
      action="New transfer"
    >
      <DataTable
        columns={['Ref', 'From', 'Technician', 'Contents', 'Status']}
        rows={checkouts}
      />
    </PageShell>
  )
}
