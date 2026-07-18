// Loading Skeleton Components
// Usage: <Skeleton.Card />, <Skeleton.Table rows={5} cols={4} />, etc.

function Base({ className = '', style }: { className?: string; style?: React.CSSProperties }) {
  return (
    <div
      className={`animate-pulse bg-gray-200 rounded ${className}`}
      style={style}
    />
  )
}

function Text({ className = '' }: { className?: string }) {
  return <Base className={`h-3 ${className || 'w-full'}`} />
}

function Title({ className = '' }: { className?: string }) {
  return <Base className={`h-5 w-48 ${className}`} />
}

function Avatar({ size = 'md' }: { size?: 'sm' | 'md' | 'lg' }) {
  const sizes = { sm: 'w-8 h-8', md: 'w-10 h-10', lg: 'w-12 h-12' }
  return <Base className={`${sizes[size]} rounded-full`} />
}

function Card({ className = '' }: { className?: string }) {
  return (
    <div className={`bg-white rounded-xl border border-gray-200 p-5 ${className}`}>
      <Base className="h-3 w-24 mb-3" />
      <Base className="h-7 w-36 mb-2" />
      <Base className="h-3 w-16" />
    </div>
  )
}

function TableRow({ cols = 4 }: { cols?: number }) {
  return (
    <div className="flex items-center gap-4 px-5 py-4 border-b border-gray-50">
      {Array.from({ length: cols }).map((_, i) => (
        <Base
          key={i}
          className={`h-4 ${i === 0 ? 'w-48' : i === cols - 1 ? 'w-20 ml-auto' : 'w-24'}`}
        />
      ))}
    </div>
  )
}

function Table({ rows = 5, cols = 4 }: { rows?: number; cols?: number }) {
  return (
    <div className="bg-white rounded-xl border border-gray-200 overflow-hidden">
      {/* Header */}
      <div className="flex items-center gap-4 px-5 py-3 border-b border-gray-100">
        {Array.from({ length: cols }).map((_, i) => (
          <Base key={i} className={`h-3 ${i === 0 ? 'w-24' : i === cols - 1 ? 'w-16 ml-auto' : 'w-16'}`} />
        ))}
      </div>
      {/* Rows */}
      {Array.from({ length: rows }).map((_, i) => (
        <TableRow key={i} cols={cols} />
      ))}
    </div>
  )
}

function BudgetCard() {
  return (
    <div className="bg-white rounded-xl border border-gray-200 p-5">
      <div className="flex items-center justify-between mb-3">
        <div className="flex items-center gap-3">
          <Base className="w-8 h-8 rounded-lg" />
          <Base className="h-4 w-32" />
        </div>
        <Base className="h-5 w-12 rounded-full" />
      </div>
      <Base className="h-2 w-full rounded-full mb-3" />
      <div className="flex items-center justify-between">
        <Base className="h-3 w-32" />
        <Base className="h-3 w-12" />
      </div>
      <Base className="h-3 w-24 mt-2" />
    </div>
  )
}

function ChartCard() {
  return (
    <div className="bg-white rounded-xl border border-gray-200 p-5">
      <Base className="h-4 w-32 mb-6" />
      <div className="flex items-end gap-3 h-40">
        {Array.from({ length: 7 }).map((_, i) => (
          <Base
            key={i}
            className="flex-1 rounded-t"
            style={{ height: `${[40, 65, 30, 80, 55, 45, 70][i]}%` }}
          />
        ))}
      </div>
    </div>
  )
}

function TxItem() {
  return (
    <div className="flex items-center justify-between px-5 py-4 border-b border-gray-50">
      <div className="flex items-center gap-4">
        <Base className="w-10 h-10 rounded-lg" />
        <div>
          <Base className="h-4 w-40 mb-1.5" />
          <Base className="h-3 w-56" />
        </div>
      </div>
      <Base className="h-4 w-24" />
    </div>
  )
}

const Skeleton = { Base, Text, Title, Avatar, Card, Table, TableRow, BudgetCard, ChartCard, TxItem }

export default Skeleton
