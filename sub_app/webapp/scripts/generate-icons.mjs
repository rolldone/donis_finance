import sharp from 'sharp'
import { readFileSync, writeFileSync } from 'fs'
import { join, dirname } from 'path'
import { fileURLToPath } from 'url'

const __dirname = dirname(fileURLToPath(import.meta.url))
const publicDir = join(__dirname, '..', 'public')

const svgBuffer = readFileSync(join(publicDir, 'favicon.svg'))

async function generateIcon(size, filename) {
  await sharp(svgBuffer)
    .resize(size, size)
    .png()
    .toFile(join(publicDir, filename))
  console.log(`✅ Generated ${filename} (${size}x${size})`)
}

async function main() {
  await generateIcon(192, 'icon-192.png')
  await generateIcon(512, 'icon-512.png')
  console.log('🎉 All icons generated!')
}

main().catch(err => { console.error(err); process.exit(1) })
