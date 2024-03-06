import * as React from 'react';
import { DataGrid } from '@mui/x-data-grid';

export default function CustomDataGrid({rows, columns, width, getRowId}) {

    if (columns == null) {
        if (rows[0] != null) {
            columns = Object.keys(rows[0]).map(key => ({
                field: key
            }));
        } else {
            columns = []
        }
    }
    
    return (
        <div style={{ width: width }}>
            {rows.length > 0 ? (
                <DataGrid
                    rows={rows}
                    columns={columns}
                    initialState={{
                    pagination: {
                        paginationModel: { page: 0, pageSize: 10 },
                    },
                    }}
                    getRowId={getRowId}
                    pageSizeOptions={[10, 25, 50, 100]}
                />
                ) : (
                    <p>No data available</p>
                )}
        </div>
    );
}