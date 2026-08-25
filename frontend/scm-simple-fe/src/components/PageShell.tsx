import type { ReactNode } from 'react'

type PageShellProps = {
  title: string
  description: string
  action?: string
  children: ReactNode
}

export function PageShell({ title, description, action, children }: PageShellProps) {
  return (
    <div className="mx-auto max-w-5xl px-8 py-8">
      <div className="mb-6 flex items-start justify-between">
        <div>
          <h1 className="text-xl font-semibold text-ink">{title}</h1>
          <p className="mt-1 text-sm text-ink-soft">{description}</p>
        </div>
        {action && (
          <button
            type="button"
            className="rounded-md  px-3.5 py-2 text-sm font-medium transition-colors hover:-soft"
          >
            {action}
          </button>
        )}
      </div>
      {children}
    </div>
  )
}

type DataTableProps = {
  columns: string[]
  rows: string[][]
}

export function DataTable({ columns, rows }: DataTableProps) {
  return (
    <div className="overflow-hidden rounded-lg border border-line bg-auto">
      <table className="w-full text-left text-sm">
        <thead>
          <tr className="border-b border-line">
            {columns.map((col) => (
              <th
                key={col}
                className="px-4 py-2.5 font-mono text-[11px] uppercase  text-ink-soft"
              >
                {col}
              </th>
            ))}
          </tr>
        </thead>
            
        <tbody>
          {rows.map((row, i) => (
            <tr key={i} className="border-b border-line last:border-0 hover:bg-auto">
              {row.map((cell, j) => (
                <td key={j} className="px-4 py-2.5 text-ink">
                  {cell}
                </td>
              ))}
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  )
}
