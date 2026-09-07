import type { Metadata } from 'next'
import './globals.css'

export const metadata: Metadata = {
  title: 'MyDrive',
  description: 'Your personal Google Drive',
}

export default function RootLayout({
  children,
}: {
  children: React.ReactNode
}) {
  return (
    <html lang="en">
      <body>{children}</body>
    </html>
  )
}
