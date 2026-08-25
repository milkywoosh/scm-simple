import { useState } from 'react'
import { Link, useLocation } from 'react-router'
import { navigation } from '../data/navigation'

export function Sidebar() {
  const location = useLocation()

  // Whichever group owns the current route starts expanded.
  const activeGroup = navigation.find((group) =>
    location.pathname.startsWith(group.basePath),
  )
  const [openGroup, setOpenGroup] = useState<string | null>(
    activeGroup?.label ?? navigation[0].label,
  )

  return (
    <aside className="flex border-spacing-2 h-full w-64 shrink-0 flex-col  ">

      <div className="flex items-center gap-2 p px-5 py-5">

        <div>
          <p className="font-mono text-sm uppercase">
            Console
          </p>
          <p className="font-mono text-sm uppercase">
            Inventory
          </p>
        </div>
      </div>

      <nav className="flex-1 overflow-y-auto px-3 py-4">
        {navigation.map((group) => {

          const isGroupActive = location.pathname.startsWith(group.basePath)
          const isOpen = openGroup === group.label

          return (
            <div key={group.label} className="mb-3">
              <div
                // type="button"
                onClick={() => setOpenGroup(isOpen ? null : group.label)}
                className={`flex-5 w-full items-center gap-3 rounded-md px-3 py-2 text-left text-sm font-medium transition-colors ${isGroupActive
                    ? '-soft -active'
                    : ' hover:-soft/60 hover:-active'
                  }`}
              >
              
                  <span className="flex-3">{group.label}</span>
                  <span
                    className={`text-xs transition-transform ${isOpen ? 'rotate-90' : ''}`}
                  >
                    &gt;
                  </span>
              </div>

              {isOpen && (
                <ul className="ml-4 mt-1 space-y-0.5 border-l border-white/10 pl-4">
                  {group.children.map((leaf) => {
                    const isLeafActive = location.pathname === leaf.path
                    return (
                      <li key={leaf.path}>
                        <Link
                          to={leaf.path}
                          className={`block rounded-md px-2.5 py-1.5 text-sm transition-colors ${isLeafActive
                              ? 'bg-signal-soft text-rail font-medium'
                              : ' hover:-active'
                            }`}
                        >
                          {leaf.label}
                        </Link>
                      </li>
                    )
                  })}
                </ul>
              )}
            </div>
          )
        })}
      </nav>
    </aside>
  )
}
