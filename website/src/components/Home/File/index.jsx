/* eslint-disable */
import get from 'lodash/get'
import noop from 'lodash/noop'
import isUndefined from 'lodash/isUndefined'
import {useContext, useState} from 'react'
import ApplicationContext from '../../../context'

import './style.css'

export default () => {
  const {
    secret,
    navigateTo,
    createSecret,
    updateSecret,
    getSecretFile,
    setAlertMessage,
  } = useContext(ApplicationContext)

  const type = 'file'
  const secretID = get(secret, 'id')
  const [note, setNote] = useState(get(secret, 'note', ''))
  const [title, setTitle] = useState(get(secret, 'title', ''))
  const [file, setFile] = useState(get(secret, 'file'), null)
  const [fileName, setFileName] = useState(get(secret, 'file_name', ''))

  const setValue = (event) => {
    const field = get(event, 'target.id', '')
    let value = get(event, 'target.value', '')

    if (field === 'file') {
      value = get(event, 'target.files[0]', null)
      const newFileName = get(event, 'target.files[0].name', '')

      setFileName(newFileName)
    }

    const changeValueHandler = {
      note  : setNote,
      title : setTitle,
      file  : setFile,
      filename  : setFileName,
    }[field] || noop

    changeValueHandler(value)
  }

  const onSaveButtonClick = async () => {
    // if secret has 'id' fieldm then we're 'updating', otherwise - 'creating'
    const apiHandler = isUndefined(secretID) ? createSecret : updateSecret

    try {
      await apiHandler({
        id : secretID,
        type,
        title,
        note,
        file,
      }, secretID)

      navigateTo('home')
    } catch (error) {
      console.warn(error.message)
      setAlertMessage(`Error: ${error.message}`)
    }
  }

  const onDownloadButtonClick = () => {
    if (fileName) {
      getSecretFile(secretID)
    }
  }

  return (
    <div className="kms-secret-edit">
      <form className="kms-secret-edit__form">
        <div className="row gy-3">
          <div className="col-md-12">
            <label htmlFor="title" className="form-label">Title</label>
            <input type="text" className="form-control" id="title" value={title} onChange={setValue} />
          </div>

          <div className="col-md-6">
            <label htmlFor="note" className="form-label">Add a few notes here</label>
            <textarea className="form-control" id="note" value={note} onChange={setValue} />
          </div>

          {secretID && fileName &&
            <div className="col-md-6  kms-secret-attachment">
              <div className="row">
                <p className="col-md-12">Attachment: {fileName}</p>
                <div className="col-md-12">
                  <button
                    type="button"
                    className="btn btn-sm btn-primary"
                    onClick={onDownloadButtonClick}
                  >
                    Download
                  </button>
                </div>
              </div>
            </div>
          }

          {!secretID &&
            <div className="col-md-6">
              <label htmlFor="note" className="form-label">Please select a file</label>
              <input
                className="form-control"
                id="file"
                type="file"
                onChange={setValue}
              />
            </div>
          }
        </div>
      </form>
      <div className="kms-secret-edit__controls">
        <button type="button" className="btn btn-primary kms-button-add" onClick={onSaveButtonClick}>
          Save
        </button>
        <button type="button" className="btn btn-danger kms-button-add" onClick={() => navigateTo('home')}>
          Cancel
        </button>
      </div>
    </div>
  )
}
