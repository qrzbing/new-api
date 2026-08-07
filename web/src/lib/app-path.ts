/*
Copyright (C) 2023-2026 QuantumNous

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as
published by the Free Software Foundation, either version 3 of the
License, or (at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program. If not, see <https://www.gnu.org/licenses/>.

For commercial licensing, please contact support@quantumnous.com
*/

function readAppBasePath(): string {
  if (typeof document === 'undefined') return ''
  return (
    document.querySelector<HTMLMetaElement>('meta[name="app-base-path"]')
      ?.content ?? ''
  )
}

export const appBasePath = readAppBasePath()

export function addAppBasePath(
  value: string,
  basePath: string = appBasePath
): string {
  if (!basePath || !value.startsWith('/') || value.startsWith('//')) {
    return value
  }
  if (
    value === basePath ||
    value.startsWith(`${basePath}/`) ||
    value.startsWith(`${basePath}?`) ||
    value.startsWith(`${basePath}#`)
  ) {
    return value
  }
  return `${basePath}${value}`
}

export function getAppBaseUrl(): string {
  if (typeof window === 'undefined') return appBasePath
  return `${window.location.origin}${appBasePath}`
}

export function resolveAppUrl(value: string): string {
  const resolvedValue = addAppBasePath(value)
  if (typeof window === 'undefined') return resolvedValue
  return new URL(resolvedValue, `${getAppBaseUrl()}/`).toString()
}
