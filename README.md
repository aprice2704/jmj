# JMJ -- a software synthesizer
### AJP Jan-2021

## Structure 

### Voices

A voice is definition for a sound, created by a sequence of Stages. Voices are not directly played, rather they are formed into Notes.

### Notes

A Note is an instance of a Voice, started at a specific point by a Trigger and shaped with a specific duration via a Sustain-Release (SR) envelope.

### Stage

Stages are the components of Voices. They may be signal sources (e.g. Sine 262.0), filters (e.g. LowPass 500.0), or sinks (e.g. Channel left) or utility (e.g. Freq MiddleC)