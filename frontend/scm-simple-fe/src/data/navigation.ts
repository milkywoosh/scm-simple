export type NavLeaf = {
  label: string
  path: string
}

export type NavGroup = {
  label: string
  code: string // two-letter rack code, used as the sidebar's signature marker
  basePath: string
  children: NavLeaf[]
}

export const navigation: NavGroup[] = [
  {
    label: 'Items',
    code: '01',
    basePath: '/items',
    children: [
      { label: 'Search', path: '/items/search' },
      { label: 'Products', path: '/items/products' },
      { label: 'Suppliers', path: '/items/suppliers' },
    ],
  },
  {
    label: 'Locations',
    code: '02',
    basePath: '/locations',
    children: [
      { label: 'Warehouse', path: '/locations/warehouse' },
      { label: 'Technician', path: '/locations/technician' },
      { label: 'Customer', path: '/locations/customer' },
    ],
  },
  {
    label: 'Transaction',
    code: '03',
    basePath: '/transfer',
    children: [
      { label: 'Warehouse', path: '/transfer/warehouse' },
      { label: 'Technician', path: '/transfer/technician' },
      { label: 'Customer', path: '/transfer/customer' },
    ],
  },
]
