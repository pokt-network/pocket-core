# Architectural Decisions Changelog

## Overview
Add an entry to this changelog to document architectural decisions related to the Pocket Core project.

## Aug 29th 2019
1.- In any `.proto` schema, account addresses will always be represented in their hex string encoding rather than a array of bytes to avoid the overhead of dealing with flexible sized bytes array.

