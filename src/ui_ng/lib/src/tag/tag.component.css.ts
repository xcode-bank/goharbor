export const TAG_STYLE = `
.sub-header-title {
  margin: 12px 0;
}

.embeded-datagrid {
  width: 98%;
  float:right; /*add for issue #2688*/
}

.hidden-tag {
  display: block; height: 0;
}

:host >>> .datagrid {
  margin: 0;
  border-bottom: none;
}

:host >>> .datagrid-placeholder {
  display: none;
}

:host >>> .datagrid .datagrid-body {
  background-color: #eee;
}

:host >>> .datagrid .datagrid-head .datagrid-row {
  background-color: #eee;
}

:host >>> .datagrid .datagrid-body .datagrid-row-master {
  background-color: #eee;
}

.truncated {
  display: inline-block;
  overflow: hidden;
  white-space: nowrap;
  text-overflow:ellipsis;
}

.copy-failed {
  color: red;
  margin-right: 6px;
}
`;