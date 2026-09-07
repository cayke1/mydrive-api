'use client'

import { useEffect, useState } from 'react'

interface HealthResponse {
  status: string
  services: {
    database: boolean
    redis: boolean
    storage: boolean
  }
}

export default function Home() {
  const [health, setHealth] = useState<HealthResponse | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    console.log('Fetching health status...')
    const fetchHealth = async () => {
      try {
        const apiUrl = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080'
        const response = await fetch(`${apiUrl}/health`, {
          method: 'GET',
          headers: {
            'Content-Type': 'application/json',
          },
        })

        if (!response.ok) {
          throw new Error(`API error: ${response.status}`)
        }

        const data = await response.json()
        setHealth(data)
      } catch (err) {
        setError(err instanceof Error ? err.message : 'Unknown error')
      } finally {
        setLoading(false)
      }
    }

    fetchHealth()
  }, [])

  return (
    <main className="min-h-screen bg-gradient-to-br from-blue-50 to-indigo-100 flex items-center justify-center p-4">
      <div className="bg-white rounded-lg shadow-lg p-8 max-w-md w-full">
        <h1 className="text-3xl font-bold text-gray-900 mb-2">MyDrive</h1>
        <p className="text-gray-600 mb-8">Your personal cloud storage</p>

        <div className="space-y-4">
          {loading && (
            <div className="text-center py-8">
              <div className="inline-block animate-spin h-8 w-8 border-4 border-indigo-500 border-t-transparent rounded-full"></div>
              <p className="text-gray-600 mt-4">Checking API connection...</p>
            </div>
          )}

          {error && (
            <div className="bg-red-50 border border-red-200 rounded-lg p-4">
              <p className="text-red-800 font-medium">Connection Error</p>
              <p className="text-red-600 text-sm mt-1">{error}</p>
            </div>
          )}

          {health && (
            <div className="space-y-4">
              <div className="flex items-center justify-between p-3 bg-gray-50 rounded-lg">
                <span className="text-gray-700 font-medium">API Status</span>
                <span
                  className={`px-3 py-1 rounded-full text-sm font-semibold ${
                    health.status === 'ok'
                      ? 'bg-green-100 text-green-800'
                      : 'bg-yellow-100 text-yellow-800'
                  }`}
                >
                  {health.status}
                </span>
              </div>

              <div className="mt-6 space-y-2">
                <p className="text-sm font-medium text-gray-700 mb-3">Service Status</p>
                <div className="grid grid-cols-1 gap-2">
                  <ServiceStatus
                    name="Database"
                    status={health.services.database}
                  />
                  <ServiceStatus
                    name="Redis Cache"
                    status={health.services.redis}
                  />
                  <ServiceStatus
                    name="Storage (MinIO)"
                    status={health.services.storage}
                  />
                </div>
              </div>
            </div>
          )}
        </div>

        <div className="mt-8 pt-6 border-t border-gray-200">
          <p className="text-xs text-gray-500 text-center">
            {new Date().toLocaleString()}
          </p>
        </div>
      </div>
    </main>
  )
}

function ServiceStatus({ name, status }: { name: string; status: boolean }) {
  return (
    <div className="flex items-center justify-between p-2 bg-gray-50 rounded">
      <span className="text-gray-700 text-sm">{name}</span>
      <span
        className={`h-2 w-2 rounded-full ${
          status ? 'bg-green-500' : 'bg-red-500'
        }`}
      ></span>
    </div>
  )
}
