import { useTranslation } from 'react-i18next'
import { PieChart, Pie, Cell, ResponsiveContainer, Tooltip, Legend } from 'recharts'
import type { CategoryBreakdown } from '../api'

const CATEGORY_COLORS = [
  '#ef4444', '#f97316', '#eab308', '#22c55e', '#06b6d4',
  '#3b82f6', '#8b5cf6', '#d946ef', '#ec4899', '#14b8a6',
  '#f43f5e', '#f59e0b', '#10b981', '#0ea5e9', '#6366f1',
  '#a855f7', '#db2777', '#78716c', '#2dd4bf', '#60a5fa',
]

interface ExpenseChartProps {
  data: CategoryBreakdown[]
}

function formatCurrency(amount: number) {
  return new Intl.NumberFormat('id-ID', {
    style: 'currency',
    currency: 'IDR',
    minimumFractionDigits: 0,
  }).format(Math.abs(amount))
}

export default function ExpenseChart({ data }: ExpenseChartProps) {
  const { t } = useTranslation()
  const expenses = data.filter((d) => d.type === 'expense' && d.total > 0)

  if (expenses.length === 0) {
    return (
      <div className="flex items-center justify-center h-48 text-sm text-muted">
        {t('common.no_data')}
      </div>
    )
  }

  const total = expenses.reduce((sum, d) => sum + d.total, 0)

  const chartData = expenses.map((d, i) => ({
    name: d.category_name || t('common.without_category'),
    value: d.total,
    color: d.category_color || CATEGORY_COLORS[i % CATEGORY_COLORS.length],
    count: d.count,
  }))

  const RADIAN = Math.PI / 180
  const renderCustomizedLabel = ({
    cx, cy, midAngle, innerRadius, outerRadius, percent,
  }: any) => {
    const radius = outerRadius + 24
    const x = cx + radius * Math.cos(-midAngle * RADIAN)
    const y = cy + radius * Math.sin(-midAngle * RADIAN)

    if (percent < 0.03) return null

    return (
      <text
        x={x}
        y={y}
        fill="var(--text-body, #374151)"
        textAnchor={x > cx ? 'start' : 'end'}
        dominantBaseline="central"
        className="text-xs"
      >
        {(percent * 100).toFixed(0)}%
      </text>
    )
  }

  return (
    <div>
      <div className="h-64">
        <ResponsiveContainer width="100%" height="100%">
          <PieChart>
            <Pie
              data={chartData}
              cx="50%"
              cy="50%"
              innerRadius={50}
              outerRadius={80}
              dataKey="value"
              labelLine={false}
              label={renderCustomizedLabel}
            >
              {chartData.map((entry, index) => (
                <Cell key={index} fill={entry.color} />
              ))}
            </Pie>
            <Tooltip
              formatter={(value: any) => formatCurrency(Number(value))}
              contentStyle={{
                backgroundColor: 'var(--bg-surface, #ffffff)',
                border: '1px solid var(--border-default, #e5e7eb)',
                borderRadius: '8px',
                fontSize: '13px',
                color: 'var(--text-body, #374151)',
              }}
            />
            <Legend
              formatter={(value: string) => (
                <span style={{ color: 'var(--text-body, #374151)', fontSize: '12px' }}>{value}</span>
              )}
            />
          </PieChart>
        </ResponsiveContainer>
      </div>
      <div className="mt-2 text-center text-sm text-muted">
        {t('common.total')}: <span className="font-semibold text-title">{formatCurrency(total)}</span>
      </div>
    </div>
  )
}
