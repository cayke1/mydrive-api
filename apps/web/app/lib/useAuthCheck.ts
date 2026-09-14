import { useRouter } from 'next/navigation'
import { useEffect } from 'react'

export function useAuthCheck() {
  const router = useRouter()

  useEffect(() => {
    // Ler o session_token do document.cookie
    const sessionToken = document.cookie
      .split('; ')
      .find(row => row.startsWith('session_token='))
      ?.split('=')[1]

    console.log(`[useAuthCheck] Session token present: ${!!sessionToken}`)

    if (!sessionToken) {
      console.log(`[useAuthCheck] No session token, redirecting to /login`)
      router.push('/login')
    }
  }, [router])
}
