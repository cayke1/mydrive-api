'use client'

import { useEffect, useRef, useState } from 'react'
import { useRouter } from 'next/navigation'
import { useAuthCheck } from '@/app/lib/useAuthCheck'
import FilePreviewModal from './FilePreviewModal'

interface File {
  id: string
  name: string
  size: number
  mime_type: string
  created_at: string
  updated_at: string
}

interface Folder {
  id: string
  name: string
  parent_id: string | null
  owner_id: string
  created_at: string
  updated_at: string
  children?: Folder[]
  files?: File[]
}

export default function FoldersPage() {
  const router = useRouter()
  useAuthCheck()

  const [rootFolders, setRootFolders] = useState<Folder[]>([])
  const [folderStack, setFolderStack] = useState<Folder[]>([])
  const [currentFolder, setCurrentFolder] = useState<Folder | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [isLoggingOut, setIsLoggingOut] = useState(false)
  const [isUploading, setIsUploading] = useState(false)
  const [previewFile, setPreviewFile] = useState<{
    meta: File
    url: string
    text?: string
  } | null>(null)
  const [isLoadingPreview, setIsLoadingPreview] = useState(false)
  const fileInputRef = useRef<HTMLInputElement>(null)
  const apiUrl = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080'

  useEffect(() => {
    fetchRootFolders()
  }, [])

  const fetchRootFolders = async () => {
    setLoading(true)
    setError(null)
    try {
      console.log('[fetchRootFolders] Starting request')
      console.log('[fetchRootFolders] API URL:', apiUrl)
      console.log('[fetchRootFolders] Document cookie:', document.cookie)

      const response = await fetch(`${apiUrl}/folders`, {
        method: 'GET',
        credentials: 'include',
        headers: { 'Content-Type': 'application/json' },
      })

      console.log('[fetchRootFolders] Response status:', response.status)
      console.log('[fetchRootFolders] Response headers:', response.headers)

      if (response.status === 401) {
        const errorText = await response.text()
        console.log('[fetchRootFolders] 401 Unauthorized. Response body:', errorText)
        router.push('/login')
        return
      }

      if (!response.ok) throw new Error(`API error: ${response.status}`)

      const data = await response.json()
      const rootOnly = (data || []).filter((f: Folder) => f.parent_id === null || !f.parent_id)
      setRootFolders(rootOnly)
      setCurrentFolder(null)
      setFolderStack([])
    } catch (err) {
      console.log('[fetchRootFolders] Error:', err)
      setError(err instanceof Error ? err.message : 'Unknown error')
    } finally {
      setLoading(false)
    }
  }

  const fetchFolderContents = async (folderId: string) => {
    setLoading(true)
    setError(null)
    try {
      const response = await fetch(`${apiUrl}/folders/${folderId}/contents`, {
        method: 'GET',
        credentials: 'include',
        headers: { 'Content-Type': 'application/json' },
      })

      if (response.status === 401) {
        router.push('/login')
        return
      }

      if (!response.ok) throw new Error(`API error: ${response.status}`)

      const data = await response.json()
      const folder: Folder = {
        ...data.folder,
        children: data.children || [],
        files: data.files || [],
      }
      setCurrentFolder(folder)
      setFolderStack([...folderStack, folder])
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Unknown error')
    } finally {
      setLoading(false)
    }
  }

  const navigateBack = () => {
    if (folderStack.length === 0) {
      setCurrentFolder(null)
    } else {
      const newStack = folderStack.slice(0, -1)
      setFolderStack(newStack)
      if (newStack.length === 0) {
        setCurrentFolder(null)
      } else {
        setCurrentFolder(newStack[newStack.length - 1])
      }
    }
  }

  const formatFileSize = (bytes: number) => {
    if (bytes === 0) return '0 B'
    const k = 1024
    const sizes = ['B', 'KB', 'MB', 'GB']
    const i = Math.floor(Math.log(bytes) / Math.log(k))
    return Math.round((bytes / Math.pow(k, i)) * 100) / 100 + ' ' + sizes[i]
  }

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString('pt-BR', {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
    })
  }

  const handleLogout = async () => {
    setIsLoggingOut(true)
    try {
      await fetch(`${apiUrl}/auth/logout`, {
        method: 'POST',
        credentials: 'include',
      })
    } finally {
      router.push('/login')
    }
  }

  const handleFileUpload = async (e: React.ChangeEvent<HTMLInputElement>) => {
    if (!currentFolder || !e.target.files?.[0]) return
    const file = e.target.files[0]
    setIsUploading(true)
    setError(null)
    const formData = new FormData()
    formData.append('folder_id', currentFolder.id)
    formData.append('file', file)
    try {
      const response = await fetch(`${apiUrl}/files/upload`, {
        method: 'POST',
        credentials: 'include',
        body: formData,
      })
      if (response.status === 401) {
        router.push('/login')
        return
      }
      if (!response.ok) throw new Error(`Upload failed: ${response.status}`)
      const newFile: File = await response.json()
      setCurrentFolder((prev) =>
        prev
          ? { ...prev, files: [...(prev.files || []), newFile] }
          : null
      )
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Upload failed')
    } finally {
      setIsUploading(false)
      if (fileInputRef.current) fileInputRef.current.value = ''
    }
  }

  const handleDownload = async (file: File) => {
    try {
      const response = await fetch(`${apiUrl}/files/${file.id}/download`, {
        method: 'GET',
        credentials: 'include',
      })
      if (response.status === 401) {
        router.push('/login')
        return
      }
      if (!response.ok) throw new Error('Download failed')
      const blob = await response.blob()
      const url = URL.createObjectURL(blob)
      const a = document.createElement('a')
      a.href = url
      a.download = file.name
      document.body.appendChild(a)
      a.click()
      document.body.removeChild(a)
      URL.revokeObjectURL(url)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Download failed')
    }
  }

  const handlePreview = async (file: File) => {
    setIsLoadingPreview(true)
    try {
      const response = await fetch(`${apiUrl}/files/${file.id}/download`, {
        credentials: 'include',
      })
      if (response.status === 401) {
        router.push('/login')
        return
      }
      if (!response.ok) throw new Error('Preview failed')
      const blob = await response.blob()
      const url = URL.createObjectURL(blob)
      let text: string | undefined
      if (file.mime_type.startsWith('text/')) text = await blob.text()
      setPreviewFile({ meta: file, url, text })
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Preview failed')
    } finally {
      setIsLoadingPreview(false)
    }
  }

  const closePreview = () => {
    if (previewFile) URL.revokeObjectURL(previewFile.url)
    setPreviewFile(null)
  }

  const handleDelete = async (file: File) => {
    if (!confirm(`Delete "${file.name}"?`)) return
    try {
      const response = await fetch(`${apiUrl}/files/${file.id}`, {
        method: 'DELETE',
        credentials: 'include',
      })
      if (response.status === 401) {
        router.push('/login')
        return
      }
      if (!response.ok) throw new Error('Delete failed')
      setCurrentFolder((prev) =>
        prev
          ? {
              ...prev,
              files: (prev.files || []).filter((f) => f.id !== file.id),
            }
          : null
      )
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Delete failed')
    }
  }

  return (
    <main className="min-h-screen bg-gray-50">
      {/* Header */}
      <div className="bg-white border-b border-gray-200 sticky top-0 z-10">
        <div className="max-w-6xl mx-auto px-8 py-4">
          <div className="flex items-center justify-between mb-4">
            <h1 className="text-3xl font-bold text-gray-900">MyDrive Explorer</h1>
            <button
              onClick={handleLogout}
              disabled={isLoggingOut}
              className="px-4 py-2 text-sm font-medium text-gray-700 bg-white border border-gray-300 rounded-lg hover:bg-gray-50 disabled:bg-gray-100 disabled:cursor-not-allowed transition-colors"
            >
              {isLoggingOut ? 'Logging out...' : 'Logout'}
            </button>
          </div>

          {/* Breadcrumb */}
          <div className="flex items-center gap-2 text-sm text-gray-600">
            <button
              onClick={fetchRootFolders}
              className="hover:text-indigo-600 hover:underline font-medium"
            >
              Root
            </button>
            {folderStack.map((folder, index) => (
              <div key={folder.id} className="flex items-center gap-2">
                <span className="text-gray-400">/</span>
                <button
                  onClick={() => {
                    const newStack = folderStack.slice(0, index + 1)
                    setFolderStack(newStack)
                    setCurrentFolder(newStack[index])
                  }}
                  className="hover:text-indigo-600 hover:underline"
                >
                  {folder.name}
                </button>
              </div>
            ))}
          </div>
        </div>
      </div>

      {/* Main Content */}
      <div className="max-w-6xl mx-auto px-8 py-8">
        {/* Controls */}
        {currentFolder && (
          <div className="mb-6 flex gap-4">
            <button
              onClick={navigateBack}
              className="inline-flex items-center gap-2 px-4 py-2 text-sm font-medium text-gray-700 bg-white border border-gray-300 rounded-lg hover:bg-gray-50 transition-colors"
            >
              <svg
                className="w-4 h-4"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M15 19l-7-7 7-7"
                />
              </svg>
              Back
            </button>
            <button
              onClick={() => fileInputRef.current?.click()}
              disabled={isUploading}
              className="inline-flex items-center gap-2 px-4 py-2 text-sm font-medium text-white bg-indigo-600 rounded-lg hover:bg-indigo-700 disabled:bg-indigo-400 transition-colors"
            >
              <svg
                className="w-4 h-4"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M9 19l3 3m0 0l3-3m-3 3v-6"
                />
              </svg>
              {isUploading ? 'Uploading...' : 'Upload File'}
            </button>
            <input
              ref={fileInputRef}
              type="file"
              onChange={handleFileUpload}
              className="hidden"
            />
          </div>
        )}

        {loading && (
          <div className="text-center py-12">
            <div className="inline-block animate-spin h-8 w-8 border-4 border-indigo-500 border-t-transparent rounded-full"></div>
            <p className="text-gray-600 mt-4">Loading...</p>
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
            {!currentFolder ? (
              // Root folders list
              <div>
                <h2 className="text-xl font-semibold text-gray-900 mb-4">
                  Root Folders
                </h2>
                <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
                  {rootFolders.length === 0 ? (
                    <div className="bg-white rounded-lg shadow p-8 text-center col-span-full">
                      <p className="text-gray-500">No folders found</p>
                    </div>
                  ) : (
                    rootFolders.map((folder) => (
                      <button
                        key={folder.id}
                        onClick={() => fetchFolderContents(folder.id)}
                        className="text-left bg-white rounded-lg shadow hover:shadow-md transition-shadow p-4 border-l-4 border-blue-500 hover:border-blue-600"
                      >
                        <div className="flex items-start justify-between mb-2">
                          <svg
                            className="w-8 h-8 text-blue-500 flex-shrink-0"
                            fill="currentColor"
                            viewBox="0 0 20 20"
                          >
                            <path d="M3 4a2 2 0 012-2h6a1 1 0 00-.82-.45l-.96-.31A1 1 0 008 2H4a2 2 0 00-2 2v12a2 2 0 002 2h12a2 2 0 002-2V4a2 2 0 00-2-2h-2.18A1 1 0 0011 2h-1a1 1 0 00-1 1v.18A2 2 0 009 2H3z" />
                          </svg>
                        </div>
                        <h4 className="font-semibold text-gray-900 truncate mb-1">
                          {folder.name}
                        </h4>
                        <p className="text-xs text-gray-500">
                          {formatDate(folder.created_at)}
                        </p>
                      </button>
                    ))
                  )}
                </div>
              </div>
            ) : (
              // Folder contents
              <div>
                <div className="mb-8">
                  <h2 className="text-2xl font-bold text-gray-900 mb-2">
                    {currentFolder.name}
                  </h2>
                  <p className="text-gray-600 text-sm">
                    {(currentFolder.children?.length || 0) +
                      (currentFolder.files?.length || 0)}{' '}
                    items
                  </p>
                </div>

                {/* Subfolders */}
                {currentFolder.children && currentFolder.children.length > 0 && (
                  <div className="mb-8">
                    <h3 className="text-lg font-semibold text-gray-800 mb-4">
                      Folders
                    </h3>
                    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
                      {currentFolder.children.map((folder) => (
                        <button
                          key={folder.id}
                          onClick={() => fetchFolderContents(folder.id)}
                          className="text-left bg-white rounded-lg shadow hover:shadow-md transition-shadow p-4 border-l-4 border-blue-500 hover:border-blue-600"
                        >
                          <div className="flex items-start justify-between mb-2">
                            <svg
                              className="w-8 h-8 text-blue-500 flex-shrink-0"
                              fill="currentColor"
                              viewBox="0 0 20 20"
                            >
                              <path d="M3 4a2 2 0 012-2h6a1 1 0 00-.82-.45l-.96-.31A1 1 0 008 2H4a2 2 0 00-2 2v12a2 2 0 002 2h12a2 2 0 002-2V4a2 2 0 00-2-2h-2.18A1 1 0 0011 2h-1a1 1 0 00-1 1v.18A2 2 0 009 2H3z" />
                            </svg>
                          </div>
                          <h4 className="font-semibold text-gray-900 truncate mb-1">
                            {folder.name}
                          </h4>
                          <p className="text-xs text-gray-500">
                            {formatDate(folder.created_at)}
                          </p>
                        </button>
                      ))}
                    </div>
                  </div>
                )}

                {/* Files */}
                {currentFolder.files && currentFolder.files.length > 0 && (
                  <div>
                    <h3 className="text-lg font-semibold text-gray-800 mb-4">
                      Files
                    </h3>
                    <div className="bg-white rounded-lg shadow overflow-hidden">
                      <table className="w-full">
                        <thead className="bg-gray-50 border-b border-gray-200">
                          <tr>
                            <th className="px-6 py-3 text-left text-xs font-medium text-gray-700 uppercase tracking-wider">
                              Name
                            </th>
                            <th className="px-6 py-3 text-left text-xs font-medium text-gray-700 uppercase tracking-wider">
                              Type
                            </th>
                            <th className="px-6 py-3 text-left text-xs font-medium text-gray-700 uppercase tracking-wider">
                              Size
                            </th>
                            <th className="px-6 py-3 text-left text-xs font-medium text-gray-700 uppercase tracking-wider">
                              Modified
                            </th>
                            <th className="px-6 py-3 text-left text-xs font-medium text-gray-700 uppercase tracking-wider">
                              Actions
                            </th>
                          </tr>
                        </thead>
                        <tbody className="divide-y divide-gray-200">
                          {currentFolder.files.map((file) => (
                            <tr
                              key={file.id}
                              className="hover:bg-gray-50 transition-colors"
                            >
                              <td className="px-6 py-4 whitespace-nowrap">
                                <div className="flex items-center gap-3">
                                  <svg
                                    className="w-5 h-5 text-gray-400"
                                    fill="currentColor"
                                    viewBox="0 0 20 20"
                                  >
                                    <path
                                      fillRule="evenodd"
                                      d="M8 16.5a.5.5 0 01-.5-.5v-5H5.707l2.147-2.146a.5.5 0 00-.708-.708l-3 3a.5.5 0 000 .708l3 3a.5.5 0 00.708-.708L5.707 11H7.5v5a.5.5 0 01-.5.5zm3-14a.5.5 0 01.5.5v5h1.793l-2.147-2.146a.5.5 0 01.708-.708l3 3a.5.5 0 010 .708l-3 3a.5.5 0 01-.708-.708L12.293 11H10.5v-5a.5.5 0 01.5-.5z"
                                      clipRule="evenodd"
                                    />
                                  </svg>
                                  <span className="text-gray-900 font-medium truncate">
                                    {file.name}
                                  </span>
                                </div>
                              </td>
                              <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                                {file.mime_type || 'Unknown'}
                              </td>
                              <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                                {formatFileSize(file.size)}
                              </td>
                              <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                                {formatDate(file.updated_at)}
                              </td>
                              <td className="px-6 py-4 whitespace-nowrap text-sm space-x-2">
                                <button
                                  onClick={() => handlePreview(file)}
                                  disabled={isLoadingPreview}
                                  className="px-3 py-1 text-xs font-medium text-indigo-600 bg-indigo-50 rounded hover:bg-indigo-100 disabled:bg-gray-100 disabled:text-gray-400 transition-colors"
                                >
                                  Preview
                                </button>
                                <button
                                  onClick={() => handleDownload(file)}
                                  className="px-3 py-1 text-xs font-medium text-blue-600 bg-blue-50 rounded hover:bg-blue-100 transition-colors"
                                >
                                  Download
                                </button>
                                <button
                                  onClick={() => handleDelete(file)}
                                  className="px-3 py-1 text-xs font-medium text-red-600 bg-red-50 rounded hover:bg-red-100 transition-colors"
                                >
                                  Delete
                                </button>
                              </td>
                            </tr>
                          ))}
                        </tbody>
                      </table>
                    </div>
                  </div>
                )}

                {/* Empty state */}
                {(!currentFolder.children || currentFolder.children.length === 0) &&
                  (!currentFolder.files || currentFolder.files.length === 0) && (
                    <div className="bg-white rounded-lg shadow p-12 text-center">
                      <svg
                        className="mx-auto h-12 w-12 text-gray-400 mb-4"
                        fill="none"
                        stroke="currentColor"
                        viewBox="0 0 24 24"
                      >
                        <path
                          strokeLinecap="round"
                          strokeLinejoin="round"
                          strokeWidth={2}
                          d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"
                        />
                      </svg>
                      <p className="text-gray-500">This folder is empty</p>
                    </div>
                  )}
              </div>
            )}
          </div>
        )}
      </div>

      {previewFile && (
        <FilePreviewModal
          file={previewFile.meta}
          url={previewFile.url}
          text={previewFile.text}
          onClose={closePreview}
        />
      )}
    </main>
  )
}
