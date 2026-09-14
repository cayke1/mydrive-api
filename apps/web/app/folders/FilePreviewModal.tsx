'use client'

import { useEffect } from 'react'

interface File {
  id: string
  name: string
  size: number
  mime_type: string
  created_at: string
  updated_at: string
}

interface FilePreviewModalProps {
  file: File
  url: string
  text?: string
  onClose: () => void
}

export default function FilePreviewModal({
  file,
  url,
  text,
  onClose,
}: FilePreviewModalProps) {
  useEffect(() => {
    const handleEscape = (e: KeyboardEvent) => {
      if (e.key === 'Escape') onClose()
    }
    document.addEventListener('keydown', handleEscape)
    return () => document.removeEventListener('keydown', handleEscape)
  }, [onClose])

  const isImage = file.mime_type.startsWith('image/')
  const isPdf = file.mime_type === 'application/pdf'
  const isText = file.mime_type.startsWith('text/')

  return (
    <div
      className="fixed inset-0 bg-black bg-opacity-50 z-50 flex items-center justify-center p-4"
      onClick={onClose}
    >
      <div
        className="bg-white rounded-lg shadow-lg max-w-4xl w-full max-h-[85vh] flex flex-col"
        onClick={(e) => e.stopPropagation()}
      >
        {/* Header */}
        <div className="flex items-center justify-between p-4 border-b border-gray-200">
          <h2 className="text-lg font-semibold text-gray-900 truncate">
            {file.name}
          </h2>
          <button
            onClick={onClose}
            className="text-gray-400 hover:text-gray-600 transition-colors"
          >
            <svg
              className="w-6 h-6"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M6 18L18 6M6 6l12 12"
              />
            </svg>
          </button>
        </div>

        {/* Content */}
        <div className="flex-1 overflow-auto p-6">
          {isImage ? (
            <div className="flex justify-center">
              <img
                src={url}
                alt={file.name}
                className="max-w-full max-h-[70vh] object-contain"
              />
            </div>
          ) : isPdf ? (
            <iframe
              src={url}
              className="w-full h-[70vh] border-0"
              title={file.name}
            />
          ) : isText && text ? (
            <pre className="overflow-auto max-h-[70vh] bg-gray-50 p-4 rounded border border-gray-200 text-sm font-mono text-gray-800 whitespace-pre-wrap break-words">
              {text}
            </pre>
          ) : (
            <div className="text-center py-12">
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
              <p className="text-gray-500">Preview not available for this file type</p>
              <p className="text-sm text-gray-400 mt-2">{file.mime_type}</p>
              <a
                href={url}
                download={file.name}
                className="mt-4 inline-block px-4 py-2 bg-indigo-600 text-white rounded-lg hover:bg-indigo-700 transition-colors text-sm font-medium"
              >
                Download File
              </a>
            </div>
          )}
        </div>
      </div>
    </div>
  )
}
