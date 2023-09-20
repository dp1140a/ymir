/**
 * Functions and Definitions for handling File Uploads
 */
import type { FilePondFile } from "filepond";

export const imageTypes: string[] = ['gif', 'jpg', 'png', 'svg'];
export const modelTypes: string[] = [
  '3ds',
  '3mf',
  'amf',
  'blend',
  'dwg',
  'dxf',
  'f3d',
  'f3z',
  'factory',
  'fcstd',
  'iges',
  'ipt',
  'obj',
  'ply',
  'py',
  'rsdoc',
  'scad',
  'shape',
  'shapr',
  'skp',
  'sldasm',
  'sldprt',
  'slvs',
  'step',
  'stl',
  'stp'
];
export const printTypes: string[] = ['gcode'];
export const otherTypes: string[] = [
  'csv',
  'doc',
  'ini',
  'json',
  'md',
  'pdf',
  'toml',
  'txt',
  'yaml',
  'yml',
  'zip'
];

export class FileUploadError extends Error {
  validExtensions: string[]
  constructor(message:string, validExtensions:string[]) {
    super(message);
    this.name = "FileUploadError";
    this.validExtensions = validExtensions
  }

}

export const CheckFileType = (fileItem:FilePondFile): FileUploadError => {
  let pondName = fileItem.getMetadata().pondName;
  let extensions: string[];
  switch (pondName) {
    case 'Image_Files':
      extensions = imageTypes;
      break;
    case 'Model_Files':
      extensions = modelTypes;
      break;
    case 'Other_Files':
      extensions = otherTypes;
      break;
    case 'Print_Files':
      extensions = printTypes;
      break;
  }
  if (extensions.includes(fileItem.fileExtension)) {
    return null;
  } else {
    let msg:string = `. ${fileItem.fileExtension} is an invalid extension for ${fileItem.getMetadata().pondName.replace('_', ' ')}`
    return new FileUploadError(msg, extensions)
  }
}