import { useLocation, useNavigate } from 'react-router'
import { navigation } from '../data/navigation'

export function Topbar() {
  const location = useLocation()
  const navigate = useNavigate()

  const group = navigation.find((g) => location.pathname.startsWith(g.basePath))
  const leaf = group?.children.find((c) => c.path === location.pathname)

  return (
    <header className="flex h-14 shrink-0 items-center justify-between border-line bg-auto px-6">
      <div className="flex items-center gap-2 text-sm">
        <button
          type="button"
          onClick={() => navigate(-1)}
          className="mr-2 rounded border border-line p px-2 py-1 font-mono text-xs text-ink-soft transition-colors hover:border-ink-soft hover:text-ink"
        >
          &larr; Back
        </button>
        
        <span className="text-ink-soft">{group?.label ?? 'Console'}</span>
        {leaf && (
          <>
            <span className="text-ink-soft/50">/</span>
            <span className="font-medium text-ink">{leaf.label}</span>
          </>
        )}
      </div>

      <div className="flex items-center gap-3">
        <span className="font-mono text-xs text-ink-soft">
          {new Date().toLocaleDateString('en-GB', {
            day: '2-digit',
            month: 'short',
            year: 'numeric',
          })}
        </span>
        <div className="h-8 w-8 rounded-full  font-mono text-xs font-medium flex items-center justify-center">
          LK
        </div>
      </div>
    </header>
  )
}
