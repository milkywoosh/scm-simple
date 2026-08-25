import { Outlet } from 'react-router'
import { Sidebar } from '../components/Sidebar'
import { Topbar } from '../components/Topbar'

export function AppLayout() {
  return (
    <div className="
      flex 
      h-screen 
      w-screen 
      overflow-hidden
      bg-auto 
      gap-x-4"
    >
      <Sidebar />

      <div className="flex min-w-2 flex-1 flex-col">
        <Topbar />
        <main className="flex-2 overflow-y-auto">
          <Outlet />
        </main>
      </div>

    </div>
  )
}
