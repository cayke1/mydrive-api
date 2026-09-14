import { NextResponse } from 'next/server'
import type { NextRequest } from 'next/server'

const AUTH_ROUTES = ['/login']

export function middleware(request: NextRequest) {
  const { pathname } = request.nextUrl
  const sessionToken = request.cookies.get('session_token')?.value

  const isAuthenticated = !!sessionToken

  // Apenas redirecionar usuários autenticados que tentam ir para /login
  if (AUTH_ROUTES.some(r => pathname.startsWith(r)) && isAuthenticated) {
    console.log(`[Middleware] Redirecting authenticated user from /login to /folders`)
    return NextResponse.redirect(new URL('/folders', request.url))
  }

  return NextResponse.next()
}

export const config = {
  matcher: ['/((?!_next/static|_next/image|favicon.ico|health).*)'],
}
