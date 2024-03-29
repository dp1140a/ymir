import type { CustomThemeConfig } from '@skeletonlabs/tw-plugin';

export const myCustomTheme: CustomThemeConfig = {
	name: 'my-custom-theme',
	properties: {
		// =~= Theme Properties =~=
		'--theme-font-family-base': `system-ui`,
		'--theme-font-family-heading': `system-ui`,
		'--theme-font-color-base': '0 0 0',
		'--theme-font-color-dark': '255 255 255',
		'--theme-rounded-base': '4px',
		'--theme-rounded-container': '4px',
		'--theme-border-base': '1px',
		// =~= Theme On-X Colors =~=
		'--on-primary': '0 0 0',
		'--on-secondary': '0 0 0',
		'--on-tertiary': '0 0 0',
		'--on-success': '0 0 0',
		'--on-warning': '0 0 0',
		'--on-error': '255 255 255',
		'--on-surface': '255 255 255',
		// =~= Theme Colors  =~=
		// primary | #7b9246
		'--color-primary-50': '235 239 227', // #ebefe3
		'--color-primary-100': '229 233 218', // #e5e9da
		'--color-primary-200': '222 228 209', // #dee4d1
		'--color-primary-300': '202 211 181', // #cad3b5
		'--color-primary-400': '163 179 126', // #a3b37e
		'--color-primary-500': '123 146 70', // #7b9246
		'--color-primary-600': '111 131 63', // #6f833f
		'--color-primary-700': '92 110 53', // #5c6e35
		'--color-primary-800': '74 88 42', // #4a582a
		'--color-primary-900': '60 72 34', // #3c4822
		// secondary | #547c99
		'--color-secondary-50': '229 235 240', // #e5ebf0
		'--color-secondary-100': '221 229 235', // #dde5eb
		'--color-secondary-200': '212 222 230', // #d4dee6
		'--color-secondary-300': '187 203 214', // #bbcbd6
		'--color-secondary-400': '135 163 184', // #87a3b8
		'--color-secondary-500': '84 124 153', // #547c99
		'--color-secondary-600': '76 112 138', // #4c708a
		'--color-secondary-700': '63 93 115', // #3f5d73
		'--color-secondary-800': '50 74 92', // #324a5c
		'--color-secondary-900': '41 61 75', // #293d4b
		// tertiary | #996699
		'--color-tertiary-50': '240 232 240', // #f0e8f0
		'--color-tertiary-100': '235 224 235', // #ebe0eb
		'--color-tertiary-200': '230 217 230', // #e6d9e6
		'--color-tertiary-300': '214 194 214', // #d6c2d6
		'--color-tertiary-400': '184 148 184', // #b894b8
		'--color-tertiary-500': '153 102 153', // #996699
		'--color-tertiary-600': '138 92 138', // #8a5c8a
		'--color-tertiary-700': '115 77 115', // #734d73
		'--color-tertiary-800': '92 61 92', // #5c3d5c
		'--color-tertiary-900': '75 50 75', // #4b324b
		// success | #7b9246
		'--color-success-50': '235 239 227', // #ebefe3
		'--color-success-100': '229 233 218', // #e5e9da
		'--color-success-200': '222 228 209', // #dee4d1
		'--color-success-300': '202 211 181', // #cad3b5
		'--color-success-400': '163 179 126', // #a3b37e
		'--color-success-500': '123 146 70', // #7b9246
		'--color-success-600': '111 131 63', // #6f833f
		'--color-success-700': '92 110 53', // #5c6e35
		'--color-success-800': '74 88 42', // #4a582a
		'--color-success-900': '60 72 34', // #3c4822
		// warning | #d3a04d
		'--color-warning-50': '248 241 228', // #f8f1e4
		'--color-warning-100': '246 236 219', // #f6ecdb
		'--color-warning-200': '244 231 211', // #f4e7d3
		'--color-warning-300': '237 217 184', // #edd9b8
		'--color-warning-400': '224 189 130', // #e0bd82
		'--color-warning-500': '211 160 77', // #d3a04d
		'--color-warning-600': '190 144 69', // #be9045
		'--color-warning-700': '158 120 58', // #9e783a
		'--color-warning-800': '127 96 46', // #7f602e
		'--color-warning-900': '103 78 38', // #674e26
		// error | #a53c23
		'--color-error-50': '242 226 222', // #f2e2de
		'--color-error-100': '237 216 211', // #edd8d3
		'--color-error-200': '233 206 200', // #e9cec8
		'--color-error-300': '219 177 167', // #dbb1a7
		'--color-error-400': '192 119 101', // #c07765
		'--color-error-500': '165 60 35', // #a53c23
		'--color-error-600': '149 54 32', // #953620
		'--color-error-700': '124 45 26', // #7c2d1a
		'--color-error-800': '99 36 21', // #632415
		'--color-error-900': '81 29 17', // #511d11
		// surface | #505050
		'--color-surface-50': '229 229 229', // #e5e5e5
		'--color-surface-100': '220 220 220', // #dcdcdc
		'--color-surface-200': '211 211 211', // #d3d3d3
		'--color-surface-300': '185 185 185', // #b9b9b9
		'--color-surface-400': '133 133 133', // #858585
		'--color-surface-500': '80 80 80', // #505050
		'--color-surface-600': '72 72 72', // #484848
		'--color-surface-700': '60 60 60', // #3c3c3c
		'--color-surface-800': '48 48 48', // #303030
		'--color-surface-900': '39 39 39' // #272727
	}
};
