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
import assert from 'node:assert/strict'
import { describe, test } from 'node:test'

import { addAppBasePath } from './app-path'

describe('application base path URL resolution', () => {
  test('prefixes app-owned root-relative paths', () => {
    assert.equal(
      addAppBasePath('/api/status', '/tools/new-api'),
      '/tools/new-api/api/status'
    )
    assert.equal(
      addAppBasePath('/oauth/discord?code=1', '/tools/new-api'),
      '/tools/new-api/oauth/discord?code=1'
    )
  })

  test('does not duplicate an existing prefix', () => {
    assert.equal(
      addAppBasePath('/new-api/api/status', '/new-api'),
      '/new-api/api/status'
    )
  })

  test('leaves root deployments and external URLs unchanged', () => {
    assert.equal(addAppBasePath('/api/status', ''), '/api/status')
    assert.equal(
      addAppBasePath('https://provider.example/callback', '/new-api'),
      'https://provider.example/callback'
    )
    assert.equal(
      addAppBasePath('//cdn.example/logo.png', '/new-api'),
      '//cdn.example/logo.png'
    )
  })
})
