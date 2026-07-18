import {
  Briefcase, Gift, Building, TrendingUp, PlusCircle,
  ShoppingCart, Truck, ShoppingBag, FileText, Music,
  Heart, Book, Home, PiggyBank, AlertTriangle,
  Folder, Wallet, Car, Plane, GraduationCap,
  Stethoscope, Gamepad2, Shirt, PawPrint, Lightbulb,
  Smartphone, Dumbbell, Baby, Wrench, Package, Landmark, CreditCard,
  type LucideProps,
} from 'lucide-react'

const ICON_MAP: Record<string, React.FC<LucideProps>> = {
  briefcase: Briefcase,
  gift: Gift,
  building: Building,
  'trending-up': TrendingUp,
  'plus-circle': PlusCircle,
  'shopping-cart': ShoppingCart,
  truck: Truck,
  'shopping-bag': ShoppingBag,
  'file-text': FileText,
  music: Music,
  heart: Heart,
  book: Book,
  home: Home,
  'piggy-bank': PiggyBank,
  'alert-triangle': AlertTriangle,
  folder: Folder,
  wallet: Wallet,
  car: Car,
  plane: Plane,
  graduation: GraduationCap,
  stethoscope: Stethoscope,
  gamepad: Gamepad2,
  shirt: Shirt,
  paw: PawPrint,
  lightbulb: Lightbulb,
  smartphone: Smartphone,
  dumbbell: Dumbbell,
  baby: Baby,
  wrench: Wrench,
  package: Package,
  landmark: Landmark,
  'credit-card': CreditCard,
}

interface CategoryIconProps {
  name?: string | null
  className?: string
}

export default function CategoryIcon({ name, className = 'w-5 h-5' }: CategoryIconProps) {
  if (!name) return <Folder className={className} />
  const Icon = ICON_MAP[name]
  if (!Icon) return <Folder className={className} />
  return <Icon className={className} />
}
