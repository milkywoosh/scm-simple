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
              <button
                type="button"
                onClick={() => setOpenGroup(isOpen ? null : group.label)}
                className={`flex w-full border border-black items-start rounded-md px-3 py-2 text-left text-sm font-medium transition-colors ${isGroupActive
                  ? '-soft -active'
                  : ' hover:-soft/60 hover:-active'
                  }`}
              >

                <span className="flex-1">{group.label}</span>

                <span
                  className={`inline-block  transition-transform duration-200 ${isOpen ? "rotate-90" : ""
                    }`}
                >

                  <svg
                    xmlns="http://www.w3.org/2000/svg"
                    fill="none"
                    viewBox="0 0 24 24"
                    strokeWidth={2}
                    stroke="currentColor"
                    className="h-4 w-4"
                  >
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      d="m8.25 4.5 7.5 7.5-7.5 7.5"
                    />
                  </svg>
                </span>
              </button>

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
