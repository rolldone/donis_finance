import { useTranslation } from 'react-i18next'
import {
  BarChart, Bar, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer, Legend,
} from 'recharts'
import type { MonthlySeriesPoint } from '../api'

interface IncomeVsExpenseChartProps {
  data: MonthlySeriesPoint[]
}

const MONTH_NAMES = ['', 'Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec']

function formatCurrency(amount: number) {
  return new Intl.NumberFormat('id-ID', {
    style: 'currency',
    currency: 'IDR',
    minimumFractionDigits: 0,
  }).format(Math.abs(amount))
}

function formatCompact(value: number): string {
  if (value >= 1_000_000) return `${(value / 1_000_000).toFixed(1)}jt`
  if (value >= 1_000) return `${(value / 1_000).toFixed(0)}rb`
  return String(value)
}

export default function IncomeVsExpenseChart({ data }: IncomeVsExpenseChartProps) {
  const { t } = useTranslation()
  const legendLabels = { income: t('transactions.income'), expense: t('transactions.expense') }

  if (!data || data.length === 0) {
    return (
      <div className="flex items-center justify-center h-48 text-sm text-muted">
        {t('common.no_data')}
      </div>
    )
  }

  const chartData = data.map((d) => ({
    label: `${MONTH_NAMES[d.month]} ${d.year}`,
    income: d.income,
    expense: d.expense,
  }))

  return (
    <div className="h-64">
      <ResponsiveContainer width="100%" height="100%">
        <BarChart data={chartData} barGap={2} barCategoryGap="20%">
          <CartesianGrid strokeDasharray="3 3" stroke="var(--border-light, #f3f4f6)" />
          <XAxis
            dataKey="label"
            tick={{ fontSize: 11, fill: 'var(--text-muted, #9ca3af)' }}
            axisLine={{ stroke: 'var(--border-light, #f3f4f6)' }}
            tickLine={false}
          />
          <YAxis
            tickFormatter={formatCompact}
            tick={{ fontSize: 11, fill: 'var(--text-muted, #9ca3af)' }}
            axisLine={false}
            tickLine={false}
          />
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
              <span style={{ color: 'var(--text-body, #374151)', fontSize: '12px' }}>{legendLabels[value as keyof typeof legendLabels] || value}</span>
            )}
          />
          <Bar dataKey="income" fill="#22c55e" radius={[4, 4, 0, 0]} maxBarSize={40} name="income" />
          <Bar dataKey="expense" fill="#ef4444" radius={[4, 4, 0, 0]} maxBarSize={40} name="expense" />
        </BarChart>
      </ResponsiveContainer>
    </div>
  )
}
