'use client'

import { useEffect, useState } from 'react'

interface Folder {
  id: string
  name: string
  created_at: string
  updated_at: string
}

export default function FoldersPage() {
  const [folders, setFolders] = useState<Folder[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    const fetchFolders = async () => {
      try {
        const apiUrl = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080'
        const response = await fetch(`${apiUrl}/folders`, {
          method: 'GET',
          headers: {
            'Content-Type': 'application/json',
          },
        })

        if (!response.ok) {
          throw new Error(`API error: ${response.status}`)
        }

        const data = await response.json()
        setFolders(data)
      } catch (err) {
        setError(err instanceof Error ? err.message : 'Unknown error')
      } finally {
        setLoading(false)
      }
    }

    fetchFolders()
  }, [])

  return (
    <main className="min-h-screen bg-gradient-to-br from-blue-50 to-indigo-100 p-8">
      <div className="max-w-4xl mx-auto">
        <h1 className="text-4xl font-bold text-gray-900 mb-2">Folders</h1>
        <p className="text-gray-600 mb-8">Your folders in MyDrive</p>

        {loading && (
          <div className="text-center py-12">
            <div className="inline-block animate-spin h-8 w-8 border-4 border-indigo-500 border-t-transparent rounded-full"></div>
            <p className="text-gray-600 mt-4">Loading folders...</p>
          </div>
        )}

        {error && (
          <div className="bg-red-50 border border-red-200 rounded-lg p-6">
            <p className="text-red-800 font-medium">Error</p>
            <p className="text-red-600 text-sm mt-1">{error}</p>
          </div>
        )}

        {!loading && !error && (
          <div>
            {folders.length === 0 ? (
              <div className="bg-white rounded-lg shadow p-8 text-center">
                <p className="text-gray-500">No folders found</p>
              </div>
            ) : (
              <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
                {folders.map((folder) => (
                  <div
                    key={folder.id}
                    className="bg-white rounded-lg shadow hover:shadow-lg transition-shadow p-6"
                  >
                    <div className="flex items-start justify-between mb-4">
                      <div className="flex-1">
                        <h3 className="text-lg font-semibold text-gray-900 truncate">
                          {folder.name}
                        </h3>
                      </div>
                      <svg
                        className="w-8 h-8 text-indigo-500 flex-shrink-0 ml-2"
                        fill="currentColor"
                        viewBox="0 0 20 20"
                      >
                        <path d="M3 4a2 2 0 012-2h6a1 1 0 00-.82-.45l-.96-.31A1 1 0 008 2H4a2 2 0 00-2 2v12a2 2 0 002 2h12a2 2 0 002-2V4a2 2 0 00-2-2h-2.18A1 1 0 0011 2h-1a1 1 0 00-1 1v.18A2 2 0 009 2H3z" />
                      </svg>
                    </div>

                    <div className="space-y-2 text-sm text-gray-600">
                      <div>
                        <p className="text-xs uppercase tracking-wide font-medium text-gray-500">
                          ID
                        </p>
                        <p className="text-gray-700 font-mono text-xs break-all">
                          {folder.id}
                        </p>
                      </div>

                      <div>
                        <p className="text-xs uppercase tracking-wide font-medium text-gray-500">
                          Created
                        </p>
                        <p className="text-gray-700">
                          {new Date(folder.created_at).toLocaleString('pt-BR')}
                        </p>
                      </div>

                      <div>
                        <p className="text-xs uppercase tracking-wide font-medium text-gray-500">
                          Updated
                        </p>
                        <p className="text-gray-700">
                          {new Date(folder.updated_at).toLocaleString('pt-BR')}
                        </p>
                      </div>
                    </div>
                  </div>
                ))}
              </div>
            )}
          </div>
        )}
      </div>
    </main>
  )
}
